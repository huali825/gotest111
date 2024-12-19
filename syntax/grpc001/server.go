package grpc001

import (
	"context"
	myGrpc "goworkwebook/syntax/grpc001/myGrpc"
)

// 保证实现了接口 保证编译能通过
var _ myGrpc.UserServiceServer = &Server{}

type Server struct {
	myGrpc.UnimplementedUserServiceServer
	Name string
}

func (s Server) GetByID(ctx context.Context, request *myGrpc.GetByIDRequest) (*myGrpc.GetByIDResponse, error) {
	//time.Sleep(time.Second)
	return &myGrpc.GetByIDResponse{
		User: &myGrpc.Person{
			Id:    123,
			Name:  "test port from: " + s.Name,
			Email: "test@test.com",
		},
	}, nil
}
