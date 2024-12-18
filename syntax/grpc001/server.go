package grpc001

import (
	"context"
	myGrpc2 "goworkwebook/syntax/grpc001/myGrpc"
)

// 保证实现了接口 保证编译能通过
var _ myGrpc2.UserServiceServer = &Server{}

type Server struct {
	myGrpc2.UnimplementedUserServiceServer
	Name string
}

func (s Server) GetByID(ctx context.Context, request *myGrpc2.GetByIDRequest) (*myGrpc2.GetByIDResponse, error) {
	return &myGrpc2.GetByIDResponse{
		User: &myGrpc2.Person{
			Id:    123,
			Name:  "test port from: " + s.Name,
			Email: "test@test.com",
		},
	}, nil
}
