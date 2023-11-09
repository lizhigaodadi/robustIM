package gateWay

import (
	"im/common/config"
	"im/gateWay/rpc/service"
	"log"
	"net"
	"time"
)

func RunMain() {
	/*TODO:Run GateWay Server*/
	config.Init("../")
	p := config.GetGateWayPort()
	log.Printf("GateWay Port : %d\n", p)
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: p,
	})
	if err != nil {
		log.Fatalf("Start TCP Epoll Server err:%s\n", err.Error())
		panic(err)
	}
	generatorInit()
	InitWPool()
	RunEpoll(ln)
	for {
		time.Sleep(10 * time.Second)
	}
}

func cmdHandler() {
	service.Init()
	cmdHandlerCount := config.GetGateWayCmdHandlerCount()

	for i := 0; i < cmdHandlerCount; i++ {
		go func() {
			for ctx := range service.Gws.Channel {
				handlerCtx(ctx)
			}
		}()
	}
}

func handlerCtx(ctx *service.CmdContext) {

}
