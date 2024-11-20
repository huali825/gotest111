package grpc001

import (
	"context"
	"goworkwebook/grpc001/myGrpc"
)

// 保证实现了接口 保证编译能通过
var _ myGrpc.UserServiceServer = &Server{}

type Server struct {
	myGrpc.UnimplementedUserServiceServer
}

func (s Server) GetByID(ctx context.Context, request *myGrpc.GetByIDRequest) (*myGrpc.GetByIDResponse, error) {
	return &myGrpc.GetByIDResponse{
		User: &myGrpc.Person{
			Id:    123,
			Name:  "test",
			Email: "test@test.com",
		},
	}, nil
}
