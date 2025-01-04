package grpcInterceptor

import (
	"context"
	"go.opentelemetry.io/otel"
	"goworkwebook/syntax/002grpcInterceptor/protobufInterface/v1"
	"log"
	"time"
)

// 保证实现了接口 保证编译能通过
var _ PtbfItfcv1.UserServiceServer = &Server{}

type Server struct {
	PtbfItfcv1.UnimplementedUserServiceServer
	Name string
}

func (s Server) GetByID(ctx context.Context, request *PtbfItfcv1.GetByIDRequest) (*PtbfItfcv1.GetByIDResponse, error) {
	//time.Sleep(time.Second)

	ctx, span := otel.Tracer("server_biz").Start(ctx, "get_by_id")
	defer span.End()

	ddl, ok := ctx.Deadline()
	if ok {
		rest := ddl.Sub(time.Now())
		log.Println(rest.String())
	}
	time.Sleep(time.Millisecond * 100)

	return &PtbfItfcv1.GetByIDResponse{
		User: &PtbfItfcv1.Person{
			Id:    123,
			Name:  "test port from: " + s.Name,
			Email: "test@test.com",
		},
	}, nil
}
