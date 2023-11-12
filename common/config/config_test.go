package config

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/spf13/viper"
	"testing"
)

func GetTest() string {
	return viper.GetString("test")
}

func TestConfig(t *testing.T) {
	err := Init()
	assert.Equal(t, err, nil)

	test := GetTest()
	fmt.Printf("test :%s\n", test)

	endpoints := GetEtcdEndPointsForDiscovery()
	fmt.Printf("len: %d\n", len(endpoints))
	for _, endpoint := range endpoints {
		fmt.Printf("endPoints: %s\n", endpoint)
	}

	timeout := GetEtcdTimeOutDialTimeForDiscovery()
	fmt.Printf("timeout: %v\n", timeout)

	ipConfigPort := GetIpConfigPort()
	assert.Equal(t, ipConfigPort, ":6789")
	fmt.Printf("ipconfigPort: %v\n", ipConfigPort)

	ipConfigPath := GetIpConfigPath()
	assert.Equal(t, ipConfigPath, "/ipConfig/dispatcher")
	fmt.Printf("ipconfigPort: %v\n", ipConfigPath)

}
