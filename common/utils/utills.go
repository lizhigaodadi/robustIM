package utils

import "im/common/config"

func GenerateIpConfigPath() string {
	return config.GetIpConfigPath() + config.GetGateWayHostPort()
}
