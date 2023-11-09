package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"im/common/prpc/discov"
	"im/logger"
	"sync"
	"sync/atomic"
)

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
	leaseId clientv3.LeaseID
	service *discov.Service
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
		s.AddService(service)

	}

}

func NewERegister(serverName string, eps []string, opt *Options) *ERegister {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:         eps,
		AutoSyncInterval:  opt.AutoSyncInterval,
		DialTimeout:       opt.dialTimeTimeOut,
		DialKeepAliveTime: opt.dialKeepAliveTimeOut,
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
