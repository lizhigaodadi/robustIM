package discovery

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"im/common/config"
	"sync"
)

type EtcdServiceRegister struct {
	cli           *clientv3.Client
	mutex         *sync.Mutex
	key           string
	val           string
	lease         clientv3.LeaseID
	keepaliveChan <-chan *clientv3.LeaseKeepAliveResponse
	ctx           *context.Context
}

func NewEtcdServiceRegister(ctx *context.Context, key string, e *EndPointInfo, lease int64) *EtcdServiceRegister {
	/*Connect To Etcd*/
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.GetEtcdEndPointsForDiscovery(),
		DialTimeout: config.GetEtcdTimeOutDialTimeForDiscovery(),
	})
	fmt.Printf("hello world\n")

	if err != nil {
		fmt.Printf("connect etcd Failed")
		return nil
	}
	meta := string(e.Marshal())
	es := &EtcdServiceRegister{
		cli:   cli,
		mutex: &sync.Mutex{},
		key:   key,
		val:   meta,
		ctx:   ctx,
	}
	/*Set Etcd Lease*/

	err = es.PutKeyWithLease(lease)

	if err != nil {
		/*Set up Lease Failed*/
		es.cli.Close()
		fmt.Printf("Connect Etcd Failed\n")
		return nil
	}
	return es
}

func (esr *EtcdServiceRegister) PutKeyWithLease(lease int64) error {
	response, err := esr.cli.Grant(*esr.ctx, lease)
	if err != nil {
		return err
	}

	/*Set atomic lease contain*/
	esr.lease = response.ID
	_, err = esr.cli.Put(*esr.ctx, esr.key, esr.val, clientv3.WithLease(response.ID))
	if err != nil {
		return nil
	}

	/*Set up auto renewal*/
	respChan, err := esr.cli.KeepAlive(*esr.ctx, esr.lease)
	esr.keepaliveChan = respChan

	return err
}

func (esr *EtcdServiceRegister) UpdateKey(info *EndPointInfo) {
	/*TODO: Modify the key value of the relevant lease*/
}

func (esr *EtcdServiceRegister) ListenToLeaseKeepalive() {
	for e := range esr.keepaliveChan {
		fmt.Printf("Lease Keepalive Event Happend Reversion: %d,LeaseId: %d\n", e.Revision, e.ID)
	}

	fmt.Printf("Lease End\n")
}

func (esr *EtcdServiceRegister) Close() error {
	/*Close Lease*/
	_, err := esr.cli.Revoke(*esr.ctx, esr.lease)
	if err != nil {
		return err
	}

	return esr.cli.Close()
}
