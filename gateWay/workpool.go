package gateWay

import (
	"fmt"
	"github.com/panjf2000/ants"
	"im/common/config"
)

var wp *ants.Pool

func InitWPool() {
	var err error
	wp, err = ants.NewPool(config.GetGateWayWorkPoolCount())
	if err != nil {
		fmt.Printf("Init WorkPool Err: %s, num:%d\n", err.Error(), config.GetGateWayWorkPoolCount())
		panic(err)
	}

}
