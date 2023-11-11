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

type ERegister struct {
	opt                Options
	serverName         string
	cli                *clientv3.Client
	perceptionServices atomic.Value
	mu                 sync.Mutex
	monitorServices    map[string]*EService
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

func RunMainRegister(ctx context.Context, service *discov.Service) {

	once.Do(func() {
		if len(service.EndPoints) == 0 {
			logger.Fatalf("Service必须有至少一个EndPoint存在，否则初始化Register模块失败")
			return
		}
		ep := service.EndPoints[0]

		if ep == nil {
			log.Fatalf("EndPoint 不能为空！！！！！")
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

		eRegister = NewERegister(serverName, discov.EndPointsToStrings(discov.NewEndPoints(ep)), NewOptions())

	})
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

	logger.Infof("Etcd Service Connection Lost ")
}

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
		return nil
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

func (e *ERegister) RunAsync() {
	go e.RunSync()
}

func (e *ERegister) RunSync() {
	go func() {
		for s := range e.registerChan {
			e.Register(s)
		}
	}()

	for s := range e.unRegisterChan {
		e.UnRegister(s)
	}

}

func (e *ERegister) Register(service *discov.Service) {
	/*TODO:将这个加入到我们管理的Service中*/
	ps := e.GetPerceptionServices()
	serviceName := service.ServiceName

	s, ok := ps[serviceName]

	if !ok { /*发现这是第一次添加该状态*/
		ps[serviceName] = service
		return
	}

	discov.AddService(s, service)

}

func (e *ERegister) UnRegister(service *discov.Service) {
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
}

func (e *ERegister) WatchService(ctx context.Context, prefixKey string) {
	watchChan := e.cli.Watch(ctx, prefixKey, clientv3.WithPrefix())

	/*开始监听变化事件*/
	for resp := range watchChan {
		for _, event := range resp.Events {
			/*TODO:开始进行反序列化*/
			service := &discov.Service{}
			err := json.Unmarshal(event.Kv.Value, service)
			if err != nil {
				logger.Warn("序列化消息发布信息失败：%v", err)
				continue
			}
			switch event.Type {
			case mvccpb.PUT:
				{
					e.RegisterService(ctx, service)
				}
			case mvccpb.DELETE:
				{
					e.UnRegisterService(ctx, service)
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

}

func (e *ERegister) NotifyListeners() {

}

func (e *ERegister) Name() string {
	return e.serverName
}

func (e *ERegister) GetService(ctx context.Context, serviceName string) *discov.Service {

	return nil
}
