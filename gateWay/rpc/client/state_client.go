package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"im/common/config"
	"im/state/rpc/service"
	"log"
)

var client service.StateClient

func ClientInit() {
	conn, err := grpc.Dial(config.GetStateServerGrpcAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Client For StateServer Init Failed\n")
		return
	}
	defer conn.Close()
	client = service.NewStateClient(conn)

}

func Cancel(ctx context.Context, endPoint string, connId int32, data []byte) (*service.StateResponse, error) {
	request := &service.StateRequest{
		Fd:       connId,
		Data:     data,
		Endpoint: endPoint,
	}

	response, err := client.CancelConn(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func SendMsg(ctx context.Context, endPoint string, connId int32, data []byte) (*service.StateResponse, error) {
	request := &service.StateRequest{
		Fd:       connId,
		Data:     data,
		Endpoint: endPoint,
	}

	response, err := client.SendMsg(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
