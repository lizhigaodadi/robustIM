package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"im/common/config"
	"im/common/prpc/discov"
	"im/logger"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var eRegister *ERegister
var once sync.Once

func GetEtcdDiscoveryInstance() discov.Discovery {
	return eRegister
}
func IsEtcdDiscoveryInit() bool {
	return eRegister != nil
}

type ERegister struct {
	opt                Options
	serverName         string
	cli                *clientv3.Client
	perceptionServices atomic.Value
	mu                 sync.Mutex           /*控制操作monitorServices的锁*/
	monitorServices    map[string]*EService /*2*/
	notifyMu           sync.RWMutex
	notifyFunc         []func()
	registerChan       chan *discov.Service
	unRegisterChan     chan *discov.Service
}

type EService struct {
	cli *clientv3.Client
	discov.Service
	serviceKey    string
	leaseId       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
}

func RunMainRegister(ctx context.Context, endPoints []string) {

	once.Do(func() {
		if len(endPoints) == 0 {
			log.Fatalf("Etcd EndPoint 不能为空")
		}
		/*开始进行初始化Register工作*/
		serverName := config.GetPrpcServerName()
		if serverName == "" {
			node, err := snowflake.NewNode(time.Now().UnixNano())
			if err != nil {
				logger.Fatalf("snowflake.NewNode 初始化失败！！！")
				return
			}
			serverName = fmt.Sprintf("RandomServerName_%s", node.Generate())

			logger.Infof("初始化一个PRPC注册服务，你应该先设置服务器的别名，我们为你随机生成了一个别名：%s", serverName)
		}

		eRegister = NewERegister(serverName, endPoints, NewOptions())

		go eRegister.handleOperatorEvent(ctx)
		go eRegister.WatchService(ctx, discov.GetPrpcPrefix())

	})
}

func (e *ERegister) handleOperatorEvent(ctx context.Context) {
	/*TODO:监听两个事件channel，当有事件发生时，判断是注册事件，还是注销事件*/
	go func() {
		for s := range e.registerChan {
			logger.Infof("收到注册事件：%v", s)
			e.RegisterMonitorService(ctx, s)
		}
	}()

	for s := range e.unRegisterChan {
		logger.Infof("收到注销事件：%v", s)
		e.UnRegisterMonitorService(ctx, s)
	}
}

func (es *ERegister) RegisterMonitorService(ctx context.Context, service *discov.Service) {
	/*空指针检查*/
	if service == nil {
		logger.Fatalf("service 不能为空")
		return
	}
	/*判断一下这是新增加还是更新操作*/
	_, ok := es.monitorServices[service.ServiceName]
	if ok {
		/*没必要了*/
		return
	}
	/*开始进行新建*/
	prpcServiceKey := discov.WithPrpcPrefix(es.serverName, service.ServiceName)

	eService := NewEService(ctx, prpcServiceKey, es.cli, service.ServiceName, service.EndPoints, es.opt.LeaseTTL)
	/*加入到我们注册的map中去*/

	es.monitorServices[service.ServiceName] = eService

}

func (es *ERegister) UnRegisterMonitorService(ctx context.Context, service *discov.Service) {
	es.mu.Lock()
	defer es.mu.Unlock()
	/*TODO:删除并停止相关的ETCD租约*/
	/*首先判断一下是否有相关的Service在里面*/
	eService, ok := es.monitorServices[service.ServiceName]
	if !ok {
		logger.Warnf("没有找到相关的Service：%s,因此这次解除租约没有必要，请检查代码", service.ServiceName)
		return
	}

	/*删除租约*/
	id := eService.leaseId
	_, err := es.cli.Revoke(ctx, id)
	if err != nil {
		logger.Fatalf("删除租约失败：%v,请检查Etcd是否在正常工作", err)
	}
	/*将相关的Service进行删除*/
	delete(es.monitorServices, service.ServiceName)

}

func (es *EService) PutKeyWithLease(ctx context.Context, key string, val string, leaseId clientv3.LeaseID) {

	_, err := es.cli.Put(ctx, key, val, clientv3.WithLease(leaseId))
	if err != nil {
		logger.Warn("put key with lease error", zap.String("key", key), zap.Error(err))
	}

	keepAlive, err := es.cli.KeepAlive(ctx, leaseId)
	if err != nil {
		logger.Warn("put key with lease error", zap.String("key", key), zap.Error(err))
		return
	}

	es.keepAliveChan = keepAlive

	go es.CheckKeepAliveChannel()

}

func (es *EService) CheckKeepAliveChannel() {
	/*TODO:开始监控心跳*/
	for resp := range es.keepAliveChan {
		logger.Infof("Receive Etcd KeepAlive Package: %v", resp)
	}

	logger.Warnf("Etcd Service Connection Lost ")
}

/*TODO(重要解释):NewEService函数会自动的将该key生成合约并发布到Etcd中去*/

func NewEService(ctx context.Context, serviceKey string, cli *clientv3.Client, serviceName string, eps []*discov.EndPoint, leaseTTL int64) *EService {
	resp, err := cli.Grant(ctx, leaseTTL)
	if err != nil {
		logger.Warn("failed to grant lease", zap.Error(err))
		return nil
	}

	leaseId := resp.ID
	service := discov.NewService(serviceName, eps)
	es := &EService{
		cli:        cli,
		Service:    *service,
		serviceKey: serviceKey,
		leaseId:    leaseId,
	}

	val, err := json.Marshal(es.Service)
	if err != nil {
		logger.Warn("failed to marshal service", zap.Error(err))
		return nil
	}

	/*TODO 发布当前这个服务到到注册机构中*/
	es.PutKeyWithLease(ctx, serviceKey, string(val), es.leaseId)
	return es
}

func (e *ERegister) GetPerceptionServices() map[string]*discov.Service {
	val := e.perceptionServices.Load()

	if val == nil {
		/*TODO:  在这里我们应该自己去新建一个map来存储进去*/
		val = make(map[string]*discov.Service)
		e.SetPerceptionService(val.(map[string]*discov.Service))
	}

	m, ok := val.(map[string]*discov.Service)
	if !ok {
		return nil
	}
	return m
}

func (e *ERegister) SetPerceptionService(m map[string]*discov.Service) {
	e.perceptionServices.Store(m)
}

func (e *ERegister) AddPerceptionServices(service *discov.Service) {
	/*TODO:判断一下该字段是否已经被初始化了*/
	ps := e.GetPerceptionServices()
	if ps == nil {
		m := make(map[string]*discov.Service)
		e.SetPerceptionService(m)
		return
	}

	/*TODO:判断一下是否已经添加过相同的服务了*/
	serviceName := service.ServiceName
	s, ok := ps[serviceName]
	if !ok { /*我们没有找到目标，直接添加即可*/
		ps[serviceName] = service
	} else { /*在原有的基础上进行更新*/
		discov.AddService(s, service)

	}

}

func NewERegister(serverName string, eps []string, opt *Options) *ERegister {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:         eps,
		AutoSyncInterval:  opt.AutoSyncInterval,
		DialTimeout:       opt.DialTimeTimeOut,
		DialKeepAliveTime: opt.DialKeepAliveTimeOut,
	})
	if err != nil {
		logger.StdLog().Fatalf("Init Etcd Register Failed")
		return nil
	}

	r := &ERegister{
		opt:                *opt,
		serverName:         serverName,
		perceptionServices: atomic.Value{},
		mu:                 sync.Mutex{},
		monitorServices:    make(map[string]*EService),
		notifyMu:           sync.RWMutex{},
		notifyFunc:         make([]func(), 0),
		registerChan:       make(chan *discov.Service),
		unRegisterChan:     make(chan *discov.Service),
		cli:                cli,
	}

	return r
}

func (e *ERegister) PushNewService(ctx context.Context, service *discov.Service) {
	/*TODO:发布一个新的Service到Etcd中供其他服务器去发现*/
	serviceKey := discov.WithPrpcPrefix(e.Name(), service.ServiceName)
	/*TODO:判断一下该服务是否已经发布过了*/

	e.mu.Lock()
	defer e.mu.Unlock()
	_, ok := e.monitorServices[service.ServiceName]
	if ok { /*We Have Push A Same Service*/
		logger.StdLog().Warnf("We Can Not Push The Same Service Twice")
		return
	}

	es := NewEService(ctx, serviceKey, e.cli, service.ServiceName, service.EndPoints, e.opt.LeaseTTL)

	/*添加进去*/
	e.monitorServices[service.ServiceName] = es
}

//func (e *ERegister) RunAsync() {
//	go e.RunSync()
//}
//
//func (e *ERegister) RunSync() {
//	go func() {
//		for s := range e.registerChan {
//			e.Register(s)
//		}
//	}()
//
//	for s := range e.unRegisterChan {
//		e.UnRegister(s)
//	}
//
//}

func (e *ERegister) Register(ctx context.Context, service *discov.Service) {
	/*TODO:将这个加入到我们管理的Service中*/
	ps := e.GetPerceptionServices()
	serviceName := service.ServiceName

	s, ok := ps[serviceName]

	if !ok { /*发现这是第一次添加该状态*/
		ps[serviceName] = service
		e.SetPerceptionService(ps)
		return
	}

	discov.AddService(s, service)
	e.SetPerceptionService(ps)

}

func (e *ERegister) UnRegister(ctx context.Context, service *discov.Service) {
	ps := e.GetPerceptionServices()
	if ps == nil {
		logger.StdLog().Warnf("PercetptionServices in %s Not Init", e.Name())
		return
	}
	serviceName := service.ServiceName
	s, ok := ps[serviceName]
	if ok {
		/*We Found The Service*/
		discov.RemoveService(s, service)
	}
	e.SetPerceptionService(ps)
}

/*TODO:！！该函数会堵塞住，请异步执行 */
func (e *ERegister) WatchService(ctx context.Context, prefixKey string) {
	watchChan := e.cli.Watch(ctx, prefixKey, clientv3.WithPrefix())

	/*开始监听变化事件*/
	for resp := range watchChan {
		for _, event := range resp.Events {
			logger.Infof("发现了新的服务事件出现")
			/*TODO:开始进行反序列化*/
			service := &discov.Service{}
			err := json.Unmarshal(event.Kv.Value, service)
			if err != nil {
				logger.Warn("序列化消息发布信息失败：%v", err)
				continue
			}
			switch event.Type {
			case mvccpb.PUT:
				{ /*TODO:这个改为加入到观测服务中去*/
					e.Register(ctx, service)
				}
			case mvccpb.DELETE:
				{
					e.UnRegister(ctx, service)
				}
			}
		}
	}
}

func (e *ERegister) RegisterService(ctx context.Context, service *discov.Service) {
	e.registerChan <- service
}

func (e *ERegister) UnRegisterService(ctx context.Context, service *discov.Service) {
	e.unRegisterChan <- service
}

func (e *ERegister) AddNotify(f func()) {
	e.notifyMu.Lock()
	defer e.notifyMu.Unlock()
	e.notifyFunc = append(e.notifyFunc, f)
}

func (e *ERegister) NotifyListeners() {
	var notifyFuncClone []func()
	e.notifyMu.RLock() /*先将需要执行的方法进行获取*/
	for _, f := range e.notifyFunc {
		notifyFuncClone = append(notifyFuncClone, f)
	}
	e.notifyMu.RUnlock()

}

func (e *ERegister) Name() string {
	return e.serverName
}

func (e *ERegister) GetService(ctx context.Context, serviceName string) *discov.Service {
	/*获取其他服务器提供的服务资源*/
	services := e.GetPerceptionServices()
	if services == nil {
		logger.Warnf("在该注册发现机构中未发现目标service：%s", serviceName)
		return nil
	}

	s, ok := services[serviceName]
	if !ok {
		logger.Warnf("在该注册发现机构中未发现目标service：%s", serviceName)
		return nil
	}

	return s
}

func (e *ERegister) ShowServiceMessage(ctx context.Context) string {
	ps := e.GetPerceptionServices()
	ret := fmt.Sprintf("{\nServerName: %s\n", e.serverName)
	ret += fmt.Sprintf("PerceptionService: {\n")
	for _, s := range ps {
		ret = fmt.Sprintf("{\nService: %s\n", s.ServiceName)
		ret += fmt.Sprintf("Service: {\n")
		for _, ep := range s.EndPoints {
			ret += fmt.Sprintf("EndPoint: {\n")
			ret += fmt.Sprintf("Ip: %s\n", ep.Ip)
			ret += fmt.Sprintf("Port: %d\n", ep.Port)
			ret += fmt.Sprintf("Weight: %d\n", ep.Weight)
			ret += fmt.Sprintf("}\n")
		}
		ret += fmt.Sprintf("}\n")
	}
	ret += fmt.Sprintf("}\n")

	ms := e.monitorServices
	ret += fmt.Sprintf("MonitorServices{\n")
	for _, s := range ms {
		ret += fmt.Sprintf("{\nService: %s\n", s.ServiceName)
		ret += fmt.Sprintf("Service: {\n")
		for _, ep := range s.EndPoints {
			ret += fmt.Sprintf("EndPoint: {\n")
			ret += fmt.Sprintf("Ip: %s\n", ep.Ip)
			ret += fmt.Sprintf("Port: %d\n", ep.Port)
			ret += fmt.Sprintf("Weight: %d\n", ep.Weight)
			ret += fmt.Sprintf("}\n")
		}
		ret += fmt.Sprintf("}\n")
	}
	ret += fmt.Sprintf("}\n")

	ret += fmt.Sprintf("}\n")

	return ret
}
