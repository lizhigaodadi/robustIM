package utils

import (
	"hash/crc32"
	"im/common/config"
)

func GenerateIpConfigPath() string {
	return config.GetIpConfigPath() + config.GetGateWayHostPort()
}

func HashStr(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func HashWithSlot(hash uint32, slot int) uint32 {
	return hash % uint32(slot)
}
