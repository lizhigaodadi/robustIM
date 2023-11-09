package gateWay

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"im/gateWay/rpc/service"
	"log"
	"net"
	"testing"
	"time"
)

type server struct {
	service.UnimplementedGatewayServer
}

func (s *server) Push(ctx context.Context, request *service.GatewayRequest) (*service.GatewayResponse, error) {
	return &service.GatewayResponse{
		Code: 200,
		Msg:  "hello world",
	}, nil
}

func (s *server) DelConn(context.Context, *service.GatewayRequest) (*service.GatewayResponse, error) {
	return &service.GatewayResponse{
		Code: 200,
		Msg:  "hello world",
	}, nil
}

func TestStateGrpcServer(t *testing.T) {
	ln, err := net.Listen("tcp", ":8972")
	if err != nil {
		log.Println("failed to listen")
		return
	}

	s := grpc.NewServer()
	service.RegisterGatewayServer(s, &server{})
	err = s.Serve(ln)
	if err != nil {
		log.Println("failed to listen")
		return
	}
	for {
		time.Sleep(10 * time.Second)
	}
}

const defaultName = "world"

var (
	addr = flag.String("addr", "127.0.0.1:8972", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func TestStateGrpcClient(t *testing.T) {
	flag.Parse()
	log.Printf("addr: %s\n", *addr)
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect \n")
	}

	defer conn.Close()
	c := service.NewGatewayClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Push(ctx, &service.GatewayRequest{
		Fd:   10,
		Data: []byte("hello world"),
	})
	if err != nil {
		log.Fatalf("did not connect \n")
		return
	}
	log.Printf("response: %s\n", r.Msg)
}
