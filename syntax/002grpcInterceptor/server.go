package grpcInterceptor

import (
	"context"
	"go.opentelemetry.io/otel"
	"log"
	"time"
)

// 保证实现了接口 保证编译能通过
var _ UserServiceServer = &Server{}

type Server struct {
	UnimplementedUserServiceServer
	Name string
}

func (s Server) GetByID(ctx context.Context, request *GetByIDRequest) (*GetByIDResponse, error) {
	//time.Sleep(time.Second)

	ctx, span := otel.Tracer("server_biz").Start(ctx, "get_by_id")
	defer span.End()

	ddl, ok := ctx.Deadline()
	if ok {
		rest := ddl.Sub(time.Now())
		log.Println(rest.String())
	}
	time.Sleep(time.Millisecond * 100)

	return &GetByIDResponse{
		User: &Person{
			Id:    123,
			Name:  "test port from: " + s.Name,
			Email: "test@test.com",
		},
	}, nil
}
