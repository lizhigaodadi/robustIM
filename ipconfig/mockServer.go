package ipconfig

import (
	"context"
	"github.com/hardcore-os/plato/common/discovery"
	"im/common/config"
	discovery2 "im/common/discovery"
	"im/common/utils"
)

var (
	//Ms *MockServer
	Ms = NewMockServer()
)

type MockServer struct {
	ips      []string
	register *discovery.ServiceRegister
}

func NewMockServer() *MockServer {
	/*TODO: MockServer Init*/
	ctx := context.Background()
	register, err := discovery.NewServiceRegister(&ctx, utils.GenerateIpConfigPath(), discovery2.EndPointInfo{
		Ip:   config.GetGateWayHost(),
		Port: config.GetGateWayPortStr(),
		Meta: make(map[string]string),
	}, 10)
	if err != nil {
		return nil
	}

	ms := &MockServer{
		ips:      []string{"192.168.5.50:8082", "192.168.5.120:8081", "192.168.5.123:8081"},
		register: register,
	}
	return ms
}

func (ms *MockServer) GetIps() []string {
	return ms.ips
}
