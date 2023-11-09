package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"im/common/prpc/discov"
	"im/logger"
	"sync"
	"sync/atomic"
)

type ERegister struct {
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

	return nil
}

func (e *ERegister) AddPerceptionServices(service *discov.Service) {
	/*TODO:判断一下该字段是否已经被初始化了*/
	val := e.perceptionServices.Load()
	if val == nil {
		/*开始进行初始化操作*/
		m := make(map[string]*discov.Service)
		e.perceptionServices.Store(m)
		val = m
	}

	/*强制类型转换*/
	ps, ok := val.(map[string]*discov.Service)
	if !ok {
		logger.StdLog().Warnf("The Perception Services Get Failed")
		return
	}

	/*TODO:判断一下是否已经添加过相同的服务了*/
}
