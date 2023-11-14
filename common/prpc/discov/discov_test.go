package discov

import (
	"context"
	"fmt"
	"github.com/magiconair/properties/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"im/common/config"
	"testing"
	"time"
)

func TestEtcdConnect(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:        []string{"127.0.0.1:2379"},
		AutoSyncInterval: 10 * time.Second,
		DialTimeout:      10 * time.Second,
	})
	assert.Equal(t, err, nil)

	response, err := cli.Get(context.Background(), "a")
	assert.Equal(t, err, nil)
	for _, kv := range response.Kvs {
		fmt.Printf("key : %s val : %s\n", string(kv.Key), string(kv.Value))
	}
}

func PreInit() {
	config.Init("../../../")
}

func TestGetPrpcName(t *testing.T) {
	PreInit()
	serverName := config.GetPrpcServerName()
	discoveryName := config.GetPrpcDiscovName()
	fmt.Printf("serverName : %s, discoveryName : %s", serverName, discoveryName)
	assert.Equal(t, serverName, "firstServer")
	assert.Equal(t, discoveryName, "etcd")

}
