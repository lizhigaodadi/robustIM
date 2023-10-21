package discovery

import (
	"context"
	"fmt"
	"github.com/magiconair/properties/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"im/common/config"
	"testing"
	"time"
)

func TestMod(t *testing.T) {
	e := &EndPointInfo{
		Key:  "hello",
		Val:  "world",
		Meta: make(map[string]string),
	}
	e.Meta["1"] = "fuck"
	e.Meta["2"] = "you"

	buf := e.Marshal()
	message := string(buf)
	fmt.Println("messsage " + message)

	newE := UnMarshal(buf)
	fmt.Printf("%v", newE)
}

func TestEtcdConnect(t *testing.T) {
	DemoInit()
	ctx := context.Background()
	register := NewEtcdServiceRegister(&ctx, "test", &EndPointInfo{
		Key:  "fuck",
		Val:  "you",
		Meta: make(map[string]string),
	}, 5)

	time.Sleep(20 * time.Second)

	register.Close()
}

func DemoInit() {
	err := config.Init()
	if err != nil {
		fmt.Printf("config init failed\n")
	}
}

func TestEtcdPutAndSet(t *testing.T) {
	DemoInit()
	/*Connect To Etcd*/
	fmt.Printf("config: %v,write: %v\n", config.GetEtcdEndPointsForDiscovery(), []string{"127.0.0.1:2379"})
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.GetEtcdEndPointsForDiscovery(),
		DialTimeout: config.GetEtcdTimeOutDialTimeForDiscovery(),
	})
	assert.Equal(t, err, nil)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = cli.Put(ctx, "key2", "value1")
	assert.Equal(t, err, nil)
	resp, err := cli.Get(ctx, "test", clientv3.WithPrefix())
	assert.Equal(t, err, nil)
	for _, kv := range resp.Kvs {
		fmt.Printf("key :%s, value: %s\n", kv.Key, kv.Value)
	}
}

func TestEtcdWatch(t *testing.T) {
	DemoInit()

	time.Sleep(5 * time.Second)

	ctx := context.Background()
	/*watcher*/
	discovery := NewEtcdServiceDiscovery(&ctx)

	go func() {
		err := discovery.WatchService("test", setTest, delTest)
		assert.Equal(t, err, nil)
	}()

	go func() { /*register*/
		ctx := context.Background()
		register := NewEtcdServiceRegister(&ctx, "test4", &EndPointInfo{
			Key:  "hello",
			Val:  "world",
			Meta: make(map[string]string),
		}, 0)

		time.Sleep(100 * time.Second)
		register.Close()
	}()

	time.Sleep(100 * time.Second)
}

func setTest(key, val string) {
	fmt.Printf("key: %s,val: %s has been set\n", key, val)
}

func delTest(key, val string) {
	fmt.Printf("key: %s,val: %s has been del\n", key, val)
}
