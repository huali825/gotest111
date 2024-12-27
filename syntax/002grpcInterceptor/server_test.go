package grpcInterceptor

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(first))
	defer func() {
		server.GracefulStop()
	}()
	userServer := &Server{}
	RegisterUserServiceServer(server, userServer)

	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	err = server.Serve(l)
	t.Log(err)
}

var first grpc.UnaryServerInterceptor = func(ctx context.Context,
	req any, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {
	log.Println("before call")
	resp, err = handler(ctx, req)
	log.Println("after call")
	return
}
