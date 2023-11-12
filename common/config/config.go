package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"net"
	"strconv"
	"time"
)

var (
	etcd_endpoints_discovery = "etcd.endpoints"
	etcd_timeout_discovery   = "etcd.timeout"
	ipconfig_port            = "ipConfig.port"
	ipconfig_path            = "ipConfig.monitorPath"
	gateway_host_port        = "gateWay.ipHost"
)

//Config File Reader

func Init() error {
	viper.AddConfigPath("../../")
	viper.SetConfigType("yaml")
	viper.SetConfigName("hucket.yaml")

	if err := viper.ReadInConfig(); err != nil {
		_, ok := err.(*viper.ConfigFileNotFoundError)
		if ok {
			return err
		} else {
			return errors.New("Parse Config File Failed ")
		}
	}

	return nil
}

/*Get Etcd EndPoints*/
func GetEtcdEndPointsForDiscovery() []string {
	endpoints := viper.GetStringSlice(etcd_endpoints_discovery)
	//fmt.Printf("---%v---\n", endpoints)
	return endpoints
}

func GetEtcdTimeOutDialTimeForDiscovery() time.Duration {
	return viper.GetDuration(etcd_timeout_discovery) * time.Second
}

/*Get IpConfig Port*/
func GetIpConfigPort() string {
	return ":" + viper.GetString(ipconfig_port)
}

func GetIpConfigPath() string {
	return viper.GetString(ipconfig_path)
}

func GetGateWayHostPort() string {
	return viper.GetString(gateway_host_port)
}
func GetGateWayHost() string {
	host, _, err := net.SplitHostPort(GetIpConfigPort())
	if err != nil {
		log.Fatalf("Parse GateWay Host Port Failed\n")
		return ""
	}

	return host
}
func GetGateWayPort() int {
	_, port, err := net.SplitHostPort(GetIpConfigPort())
	if err != nil {
		log.Fatalf("Parse GateWay Host Port Failed\n")
		return 0
	}

	p, err := strconv.ParseInt(port, 10, 16)
	if err != nil {
		log.Fatalf("Parse GateWay Host Port Failed\n")
		return 0
	}

	return int(p)
}

func GetGateWayPortStr() string {
	_, port, err := net.SplitHostPort(GetIpConfigPort())
	if err != nil {
		log.Fatalf("Parse GateWay Host Port Failed\n")
		return "0"
	}
	return port
}
