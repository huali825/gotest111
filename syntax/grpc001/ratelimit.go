package grpc001

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

type TokenBucketLimiter struct { // 令牌桶
	interval  time.Duration // 间隔时间
	buckets   chan struct{}
	closeCh   chan struct{}
	closeOnce sync.Once
}

// NewSvcInterceptor 创建一个新的服务拦截器
func (c *TokenBucketLimiter) NewSvcInterceptor() grpc.UnaryServerInterceptor {
	// 创建一个定时器，每隔c.interval时间就向c.buckets通道发送一个空结构体
	ticker := time.NewTicker(c.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				select {
				case c.buckets <- struct{}{}:
					// 如果c.buckets通道未满，则向其发送一个空结构体
				default:
					// 如果c.buckets通道已满，则不做任何操作
				}
			case <-c.closeCh:
				// 如果接收到关闭信号，则退出循环
				return
			}
		}
	}()

	// 返回一个拦截器函数
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		select {
		case <-c.buckets:
			// 如果c.buckets通道中有空结构体，则调用handler函数处理请求
			return handler(ctx, req)
		//做法1
		default:
			// 如果c.buckets通道已满，则返回一个ResourceExhausted错误
			return nil, status.Errorf(codes.ResourceExhausted, "限流")
			// 做法2
			// 如果请求上下文已关闭，则返回一个上下文错误
			//case <-ctx.Done():
			//	return nil, ctx.Err()
		}
	}
}
