package logging

import (
	"context"
	"google.golang.org/grpc"
)

type InterceptorBuilder struct {
	name string
}

func (i *InterceptorBuilder) Build() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		return handler(ctx, req)
	}
}
