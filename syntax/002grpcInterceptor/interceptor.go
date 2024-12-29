package grpcInterceptor

/*
 * @Date: 2024年12月27日
 * @LastEditors: TMH
 * @LastEditTime: 2024年12月29日12:51:27
 * @FilePath: syntax/002grpcInterceptor/interceptor.go
 * @Description: grpc.UnaryServerInterceptor Server 拦截器  针对特定方法限流, 耦合性强 不通用
 */

import (
	"context"
	"fmt"
	"goworkwebook/webook003/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"goworkwebook/webook003/pkg/ratelimit"
)

type InterceptorBuilder struct {
	limiter ratelimit.Limiter
	key     string
	l       logger.LoggerV1
}

func NewInterceptorBuilder(limiter ratelimit.Limiter, key string, l logger.LoggerV1) *InterceptorBuilder {
	return &InterceptorBuilder{limiter: limiter, key: key, l: l}
}

// BuildServerInterceptorServiceBiz 针对业务限流
func (b *InterceptorBuilder) BuildServerInterceptorServiceBiz() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

		if idReq, ok := req.(*GetByIDRequest); ok {
			limited, err := b.limiter.Limit(ctx,
				// limiter:user:456
				fmt.Sprintf("limiter:user:%s:%d", info.FullMethod, idReq.Id))
			if err != nil {
				// err 不为nil，你要考虑你用保守的，还是用激进的策略
				// 这是保守的策略
				b.l.Error("判定限流出现问题", logger.Error(err))
				return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
				// 这是激进的策略
				// return handler(ctx, req)
			}
			if limited {
				return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
			}
		}
		return handler(ctx, req)
	}
}
