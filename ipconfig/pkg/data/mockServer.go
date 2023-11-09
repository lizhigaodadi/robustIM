package data

import (
	"context"
	"im/common/config"
	discovery "im/common/discovery"
	"log"
	"math/rand"
	"net"
	"time"
)

type MockServer struct {
	ips      []string
	register *discovery.EtcdServiceRegister
	ip       string
	port     string
}

func NewMockServer(hostPort string) *MockServer {
	/*TODO: MockServer Init*/
	ctx := context.Background()

	host, p, err := net.SplitHostPort(hostPort)
	if err != nil {
		log.Printf("MockServer Init Failed: Parse HostPort Failed\n")
		return nil
	}

	register := discovery.NewEtcdServiceRegister(&ctx, "/ipConfig/dispatcher"+hostPort, &discovery.EndPointInfo{
		Ip:   config.GetGateWayHost(),
		Port: config.GetGateWayPortStr(),
		Meta: make(map[string]interface{}),
	}, 10)
	if register == nil {
		log.Fatalf("Etcd Register Service Init Failed\n")
		return nil
	}

	ms := &MockServer{
		ips:      []string{"192.168.5.50:8082", "192.168.5.120:8081", "192.168.5.123:8081"},
		register: register,
		ip:       host,
		port:     p,
	}
	log.Printf("Init MockServer Successfully\n")
	return ms
}

func (ms *MockServer) GetIps() []string {
	log.Printf("MockServer GetIps\n")
	return ms.ips
}

func (ms *MockServer) UpdateData() {
	for {
		pointInfo := &discovery.EndPointInfo{
			Ip:   ms.ip,
			Port: ms.port,
			Meta: make(map[string]interface{}),
		}
		pointInfo.Meta["messageBytes"] = float64(rand.Int63n(123123112123))
		pointInfo.Meta["connNum"] = float64(rand.Int63n(11111123))

		ms.register.UpdateKey(pointInfo)

		log.Printf("Register Update Data: %v\n", pointInfo)

		time.Sleep(2 * time.Second)
	}
}
