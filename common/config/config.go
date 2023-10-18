package config

import "github.com/spf13/viper"

//Config File Reader

func Init(path string) {
	viper.AddConfigPath("../../")
	viper.SetConfigType("yaml")
	viper.SetConfigName("plato.yaml")

	if err := viper.ReadInConfig(); err != nil {

	}

}
