package ipconfig

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"im/common/config"
	"im/ipconfig/pkg/data"
	"testing"
	"time"
)

func TestIpConfig(t *testing.T) {

	Init()
	RegisterService()

	time.Sleep(100 * time.Second)

}

func TestConfig(t *testing.T) {
	err := config.Init("../")
	assert.Equal(t, err, nil)

	endpoints := config.GetEtcdEndPointsForDiscovery()
	fmt.Printf("len: %d\n", len(endpoints))
	for _, endpoint := range endpoints {
		fmt.Printf("endPoints: %s\n", endpoint)
	}

	timeout := config.GetEtcdTimeOutDialTimeForDiscovery()
	fmt.Printf("timeout: %v\n", timeout)

	ipConfigPort := config.GetIpConfigPort()
	assert.Equal(t, ipConfigPort, ":6789")
	fmt.Printf("ipconfigPort: %v\n", ipConfigPort)

	ipConfigPath := config.GetIpConfigPath()
	assert.Equal(t, ipConfigPath, "/ipConfig/dispatcher")
	fmt.Printf("ipconfigPort: %v\n", ipConfigPath)

}

func RegisterService() {
	server1 := data.NewMockServer("192.168.0.1:92")
	go server1.UpdateData()
	server2 := data.NewMockServer("192.168.0.2:92")
	go server2.UpdateData()
	server3 := data.NewMockServer("192.168.0.3:92")
	go server3.UpdateData()
	server4 := data.NewMockServer("192.168.0.4:92")
	go server4.UpdateData()
	server5 := data.NewMockServer("192.168.0.5:92")
	go server5.UpdateData()
}
