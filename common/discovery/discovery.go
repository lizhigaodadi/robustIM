package discovery

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"im/common/config"
	"sync"
)

type EtcdServiceDiscovery struct {
	cli   *clientv3.Client
	mutex *sync.Mutex
	ctx   *context.Context
}

func NewEtcdServiceDiscovery(ctx *context.Context) *EtcdServiceDiscovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.GetEtcdEndPointsForDiscovery(),
		DialTimeout: config.GetEtcdTimeOutDialTimeForDiscovery(),
	})
	if err != nil {
		/*TODO: handle error*/
		return nil
	}

	return &EtcdServiceDiscovery{
		ctx:   ctx,
		mutex: &sync.Mutex{},
		cli:   cli,
	}
}

func (esd *EtcdServiceDiscovery) WatchService(prefix string, set, del func(key, val string)) error {
	/*TODO: Add new events to be monitored*/

	resp, err := esd.cli.Get(*esd.ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		/*TODO: err handler*/
		return err
	}

	for _, kv := range resp.Kvs {
		set(string(kv.Key), string(kv.Value))
	}

	/*Sign up for related watch events*/
	esd.watcher(prefix, resp.Header.Revision+1, set, del)
	return nil
}

func (s *EtcdServiceDiscovery) watcher(prefix string, rev int64, set, del func(key, val string)) {
	rch := s.cli.Watch(*s.ctx, prefix, clientv3.WithPrefix(), clientv3.WithRev(rev))

	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				set(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				del(string(ev.Kv.Key), string(ev.Kv.Value))
			}
		}
	}

}

func (s *EtcdServiceDiscovery) Close() error {
	/*Turn off listening event*/
	return s.cli.Close()
}
