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
	gateway_conn_limit       = "gateWay.connLimit"
	gateway_queue_size       = "gateWay.queueSize"
	gateway_epoll_count      = "gateWay.epollCount"
	gateway_port             = "gateWay.port"
	gateway_work_pool_count  = "gateWay.workPoolCount"
	gateway_handler_count    = "gateWay.cmdHandlerCount"
	gateway_grpc_addr        = "gateWay.grpc.address"
	state_grpc_addr          = "stateServer.grpc.address"
	state_login_slot         = "stateServer.loginSlot"
	redis_endpoints          = "redis.endpoints"
	prpc_discov_name         = "etcd"
	prpc_server_name         = "prpc.serverName"
)

//Config File Reader

func Init(path string) error {
	viper.AddConfigPath(path)
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

	log.Printf("Config Module Init Successfully\n")

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

func GetGateConnLimit() int32 {
	return viper.GetInt32(gateway_conn_limit)
}

func GetGateWayQueueSize() int {
	return viper.GetInt(gateway_queue_size)
}

func GetGateWayEpollCount() int {
	return viper.GetInt(gateway_epoll_count)
}

func GetGateWayWorkPoolCount() int {
	return viper.GetInt(gateway_work_pool_count)
}

func GetGateWayCmdHandlerCount() int {
	return viper.GetInt(gateway_handler_count)
}

func GetStateServerGrpcAddr() string {
	return "addr: " + viper.GetString(state_grpc_addr)
}

func GetRedisEndpoint() []string {
	return viper.GetStringSlice(redis_endpoints)
}

func GetGateWayGrpcAddr() string {
	return viper.GetString(gateway_grpc_addr)
}

func GetStateServerLoginSlot() int {
	return viper.GetInt(state_login_slot)
}

func GetPrpcServerName() string {
	return viper.GetString(prpc_server_name)
}

/*获取当前是选择了什么作为注册机构*/

func GetPrpcDiscovName() string {
	return viper.GetString(prpc_discov_name)
}
