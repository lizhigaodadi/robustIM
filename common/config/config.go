package config

import (
	"errors"
	"github.com/spf13/viper"
	"time"
)

var (
	etcd_endpoints_discovery = "etcd.endpoints"
	etcd_timeout_discovery   = "etcd.timeout"
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
