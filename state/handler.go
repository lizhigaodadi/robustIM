package state

import (
	"github.com/golang/protobuf/proto"
	"im/idl/message"
	"im/state/relation"
	"im/state/rpc/service"
	"log"
)

const (
	AckSuccess = 200
	AckErr     = 300
)

const (
	RedisLoginKey = "%d_login_%d"
)

func contextHandler(ctx *service.CmdContext) {
	/* handler different Event*/
	switch ctx.CmdType {
	case service.AckType:
		{
			AckHandler(ctx)
		}
	case service.ReConnType:
		{
			ReConnHandler(ctx)
		}
	case service.LoginOutType:
		{
			LoginOutHandler(ctx)
		}
	case service.HeartBeatType:
		{
			HeartBeatHandler(ctx)
		}
	case service.PushType:
		{
			PushHandler(ctx)
		}
	case service.UpType:
		{
			UpHandler(ctx)
		}
	}
}

func LoginHandler(ctx *service.CmdContext) {
	/*TODO: Login Handler*/
	/*TODO:UnMarshal Login Message*/
	if ctx.CmdType != message.CmdType_Login {
		log.Printf("Not Matched Message Type, TargetType: %s,ActualType: %s\n", message.CmdType_name[int32(message.CmdType_Login)], message.CmdType_name[int32(ctx.CmdType)])
	}
	/*Create An New Connect State*/
	state := NewConnState(ctx.ConnId)
	cacheState.AddConn(state)
	connId := ctx.ConnId
	relation.UpdateLogin(ctx.Ctx)

}

func AckHandler(ctx *service.CmdContext) {
	connId := ctx.ConnId
	if ctx.CmdType != message.CmdType_ACK {
		panic("Handler Handle Fatal Message")
	}

	body := ctx.Data
	// UMarshal
	AckMessage := &message.ACKMsg{}
	err := proto.Unmarshal(body, AckMessage)
	if err != nil {
		log.Printf("log ")
	}

	state := cacheState.GetConnState(connId)
	if state == nil {
		log.Printf("ConnId: %d Has Been Disabled", connId)
		/*TODO: MayBe We Can Told Client THis ConnId is disabled?*/
		return
	}

	if AckMessage.Code == AckSuccess {
		/*Successfully!*/
		state.CloseHeartbeatTimer()
	} else {
		/*TODO: We does not has an good idea about this module*/
	}

}

func ReConnHandler(ctx *service.CmdContext) {
	connId := ctx.ConnId /*This Is An Old Conn*/
	if ctx.CmdType != message.CmdType_ReConn {
		panic("Type Not Match!")
	}

	/*UnMarshal*/
	reConnMsg := &message.ReConnMsg{}
	err := proto.Unmarshal(ctx.Data, reConnMsg)
	if err != nil {
		log.Printf("UnMarshal ReConnMsg Failed\n")
		return
	}
	/*Check This Conn Still Contains?*/
	state := cacheState.GetConnState(connId)

	if state == nil {
		/*This Conn State has been Deleted*/
		/*Create An New ConnState*/
		connState := NewConnState(connId)
		cacheState.AddConn(connState)

		/*Update Login Status*/

	} else {
		/*Old Conn State Still Contains */
		state.CloseReConnTimer()
		/*TODO:Send Ack Package*/
	}
}

func LoginOutHandler(ctx *service.CmdContext) {

}

func HeartBeatHandler(ctx *service.CmdContext) {

}

func PushHandler(ctx *service.CmdContext) {

}

func UpHandler(ctx *service.CmdContext) {

}
