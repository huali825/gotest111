package grpc001

import (
	"google.golang.org/grpc"
	"goworkwebook/syntax/grpc001/myGrpc"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	server := grpc.NewServer()
	userServer := &Server{}
	myGrpc.RegisterUserServiceServer(server, userServer)

	l, err := net.Listen("tcp", "127.0.0.1:8090")
	if err != nil {
		t.Fatal(err)
	}
	_ = server.Serve(l)
}
