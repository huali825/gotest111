package gozero

import (
	"context"
)

// 保证实现了接口 保证编译能通过
var _ UserServiceServer = &Server{}

type Server struct {
	UnimplementedUserServiceServer
	Name string
}

func (s Server) GetByID(ctx context.Context, request *GetByIDRequest) (*GetByIDResponse, error) {
	return &GetByIDResponse{
		User: &Person{
			Id:    123,
			Name:  "test",
			Email: "test@test.com",
		},
	}, nil
}
