package service

import (
	"context"
	"github.com/hardcore-os/plato/common/idl/message"
	"im/common/config"
)

const (
	DEL_TYPE  = 0
	PUSH_TYPE = 1
)

var Gws *GateWayServerService

type CmdContext struct {
	ctx    context.Context
	connId uint32
	ct     message.CmdType
	data   []byte
}

type GateWayServerService struct {
	Channel chan *CmdContext
}

func Init() {
	Gws = &GateWayServerService{
		Channel: make(chan *CmdContext, config.GetGateWayCmdHandlerCount()),
	}
}

func (gss *GateWayServerService) DelConn(ctx context.Context, request *GatewayRequest) (*GatewayResponse, error) {

	gss.Channel <- &CmdContext{
		ctx:    ctx,
		ct:     DEL_TYPE,
		connId: uint32(request.Fd),
		data:   request.Data,
	}

	return &GatewayResponse{
		Code: 200,
		Msg:  "success",
	}, nil
}

func (gss *GateWayServerService) Push(ctx context.Context, request *GatewayRequest) (*GatewayResponse, error) {
	gss.Channel <- &CmdContext{
		ctx:    ctx,
		ct:     PUSH_TYPE,
		connId: uint32(request.Fd),
		data:   request.Data,
	}

	return &GatewayResponse{
		Code: 200,
		Msg:  "success",
	}, nil
}
