package service

import (
	"context"
	"im/idl/message"
)

const (
	CANCEL_TYPE = 0
	SEND_TYPE   = 1
)

/*CmdType*/
const (
	LoginType     = 0
	LoginOutType  = 1
	ReConnType    = 2
	HeartBeatType = 3
	AckType       = 4
	PushType      = 5
	UpType        = 6
)

var Sss *StateServerService

type CmdContext struct {
	Ctx     context.Context
	ConnId  uint64
	CmdType message.CmdType
	Data    []byte
}

type StateServerService struct {
	Channel chan *CmdContext
}

func ServiceInit() {
	Sss = &StateServerService{
		Channel: make(chan *CmdContext),
	}
}

func (sss *StateServerService) CancelConn(ctx context.Context, request *StateRequest) (*StateResponse, error) {
	connId := request.Fd
	cc := &CmdContext{
		Ctx:     ctx,
		ConnId:  uint64(connId),
		CmdType: CANCEL_TYPE,
		Data:    request.Data,
	}
	sss.Channel <- cc
	return &StateResponse{
		Code: 200,
		Msg:  "success",
	}, nil
}

func (sss *StateServerService) SendMsg(ctx context.Context, request *StateRequest) (*StateResponse, error) {
	connId := request.Fd
	cc := &CmdContext{
		Ctx:     ctx,
		ConnId:  uint64(connId),
		CmdType: SEND_TYPE,
		Data:    request.Data,
	}

	sss.Channel <- cc

	return &StateResponse{
		Code: 200,
		Msg:  "success",
	}, nil
}
