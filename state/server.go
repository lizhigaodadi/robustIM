package state

import (
	"im/state/rpc/service"
)

func RunMain() {
	service.ServiceInit()
	CmdChannelListening()
}

func CmdChannelListening() {
	go func() {
		for context := range service.Sss.Channel {
			contextHandler(context)
		}
	}()
}
