package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"im/common/config"
	"im/gateWay/rpc/service"
	"log"
)

var client service.GatewayClient

func ClientInit() {
	conn, err := grpc.Dial(config.GetGateWayGrpcAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Init GateWayGrpcClient Failed\n")
		return
	}
	defer conn.Close()

	client = service.NewGatewayClient(conn)
}

/*Remote Process Call*/

func DelConn(ctx context.Context, request *service.GatewayRequest) (*service.GatewayResponse, error) {
	response, err := client.DelConn(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, err
}

func Push(ctx context.Context, request *service.GatewayRequest) (*service.GatewayResponse, error) {
	response, err := client.Push(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, err
}
