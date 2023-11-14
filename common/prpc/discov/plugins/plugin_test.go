package plugins

import (
	"context"
	"fmt"
	"github.com/magiconair/properties/assert"
	"im/common/config"
	discov2 "im/common/prpc/discov"
	"testing"
	"time"
)

func PreInit() {
	config.Init("../../../../")
}

func TestEtcdRegisterAndDiscovery(t *testing.T) {
	/*测试Prpc的Etcd服务发现注册的可用性*/
	PreInit()
	discov, err := GetDiscoInstance(context.Background())
	discov.RegisterService(context.Background(), &discov2.Service{
		ServiceName: "test",
		EndPoints: []*discov2.EndPoint{&discov2.EndPoint{
			Port:       8080,
			Ip:         "192.168.5.9",
			Weight:     100,
			ServerName: "王小美的Windows",
		}},
	})

	/*开始进行服务注册*/
	go func() {
		discov, err := GetDiscoInstance(context.Background())
		assert.Equal(t, err, nil)
		s := &discov2.Service{
			ServiceName: "test2",
			EndPoints: []*discov2.EndPoint{&discov2.EndPoint{
				Port:       8080,
				Ip:         "192.168.5.10",
				Weight:     100,
				ServerName: "王小美的MacOS",
			}}}

		discov.RegisterService(context.Background(), s)
		time.Sleep(time.Second * 20)
		discov.UnRegisterService(context.Background(), s)
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			message := discov.ShowServiceMessage(context.Background())
			fmt.Println(message)
		}
	}()

	assert.Equal(t, err, nil)
	time.Sleep(100 * time.Second)

}

func TestEtcdRegister(t *testing.T) {
	PreInit()
	discov, err := GetDiscoInstance(context.Background())
	discov.RegisterService(context.Background(), &discov2.Service{
		ServiceName: "test3",
		EndPoints: []*discov2.EndPoint{&discov2.EndPoint{
			Port:       8080,
			Ip:         "192.168.5.11",
			Weight:     100,
			ServerName: "王小美的另一台Windows",
		}},
	})
	go func() {
		for {
			time.Sleep(5 * time.Second)
			message := discov.ShowServiceMessage(context.Background())
			fmt.Println(message)
		}
	}()

	assert.Equal(t, err, nil)
	time.Sleep(1000 * time.Second)

}
