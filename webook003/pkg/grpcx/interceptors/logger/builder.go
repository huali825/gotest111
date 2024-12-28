package logger

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"goworkwebook/webook003/pkg/logger"
	"runtime"
	"time"
)

// InterceptorBuilder 是一个结构体，用于构建拦截器
type InterceptorBuilder struct {
	l logger.LoggerV1 // l 是一个日志记录器，遵循 LoggerV1 接口
	//interceptors.Builder                 // 嵌入 Builder 接口，提供拦截器构建功能
}

func NewInterceptorBuilder(l logger.LoggerV1) *InterceptorBuilder {
	return &InterceptorBuilder{l: l}
}

// BuildServerUnaryInterceptor 构建一个服务端的一元拦截器
func (b *InterceptorBuilder) BuildServerUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now() // 记录当前时间，用于计算请求处理时间
		event := "normal"   // 初始化事件类型为正常

		defer func() {
			// 最终输出日志
			cost := time.Since(start) // 计算请求处理时间

			// 发生了 panic
			if rec := recover(); rec != nil {
				switch re := rec.(type) {
				case error:
					err = re // 如果 panic 是 error 类型，直接赋值给 err
				default:
					err = fmt.Errorf("%v", rec) // 否则，将 panic 转换为 error
				}
				event = "recover"                                                 // 将事件类型设置为 recover
				stack := make([]byte, 4096)                                       // 创建一个字节切片用于存储堆栈信息
				stack = stack[:runtime.Stack(stack, true)]                        // 获取当前堆栈信息
				err = status.New(codes.Internal, "panic, err "+err.Error()).Err() // 创建一个新的状态错误
			}

			fields := []logger.Field{
				// unary stream 是 grpc 的两种调用形态
				logger.String("type", "unary"),            // 记录调用类型为 unary
				logger.Int64("cost", cost.Milliseconds()), // 记录请求处理时间（毫秒）
				logger.String("event", event),             // 记录事件类型
				logger.String("method", info.FullMethod),  // 记录调用的方法名
				// 客户端的信息 需要客户端配合
				//需要知道是客户端哪个业务调用过来的
				//logger.String("peer", b.PeerName(ctx)),  // 记录客户端名称
				//logger.String("peer_ip", b.PeerIP(ctx)), // 记录客户端 IP
			}
			st, _ := status.FromError(err) // 从错误中获取状态
			if st != nil {
				// 错误码
				fields = append(fields, logger.String("code", st.Code().String())) // 记录错误码
				fields = append(fields, logger.String("code_msg", st.Message()))   // 记录错误信息
			}

			b.l.Info("RPC调用", fields...) // 输出日志
		}()
		resp, err = handler(ctx, req) // 调用实际的处理器
		return
	}
}
