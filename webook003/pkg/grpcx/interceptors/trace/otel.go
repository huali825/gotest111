package trace

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"goworkwebook/webook003/pkg/grpcx/interceptors"
)

// OTELInterceptorBuilder 结构体定义了一个 OpenTelemetry 拦截器构建器，用于创建 OpenTelemetry 拦截器
type OTELInterceptorBuilder struct {
	// tracer 用于创建 OpenTelemetry 跟踪器
	tracer trace.Tracer

	// propagator 用于创建 OpenTelemetry 传播器
	propagator propagation.TextMapPropagator

	// Builder 用于创建 OpenTelemetry 拦截器
	interceptors.Builder

	// serviceName 用于指定 OpenTelemetry 拦截器的服务名称
	serviceName string
}

// 创建一个新的OTELInterceptorBuilder实例
func NewOTELInterceptorBuilder(
	// 服务名称
	serviceName string,
	// 跟踪器
	tracer trace.Tracer,
	// 传播器
	propagator propagation.TextMapPropagator) *OTELInterceptorBuilder {
	// 返回一个新的OTELInterceptorBuilder实例
	return &OTELInterceptorBuilder{tracer: tracer,
		// 设置服务名称
		serviceName: serviceName,
		// 设置传播器
		propagator: propagator}
}

func (b *OTELInterceptorBuilder) BuildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	// 获取tracer
	tracer := b.tracer
	// 如果tracer为空，则使用otel.Tracer创建一个新的tracer
	if tracer == nil {
		tracer = otel.Tracer("gitee.com/geekbang/basic-go/webook/pkg/grpcx")
	}
	// 获取propagator
	propagator := b.propagator
	// 如果propagator为空，则使用otel.GetTextMapPropagator创建一个新的propagator
	if propagator == nil {
		propagator = otel.GetTextMapPropagator()
	}
	// 创建一组属性
	attrs := []attribute.KeyValue{
		// 设置RPC系统为grpc
		semconv.RPCSystemKey.String("grpc"),
		// 设置RPC类型为unary
		attribute.Key("rpc.grpc.kind").String("unary"),
		// 设置RPC组件为server
		attribute.Key("rpc.component").String("server"),
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply interface{}, err error) {
		// 从上下文中提取信息
		ctx = extract(ctx, propagator)
		// 创建一个新的span，并设置span的名称为info.FullMethod
		ctx, span := tracer.Start(ctx, info.FullMethod,
			// 设置span的属性
			trace.WithAttributes(attrs...),
			// 设置span的类型为server
			trace.WithSpanKind(trace.SpanKindServer))
		// 在函数结束时结束span
		defer func() {
			span.End()
		}()
		// 设置span的属性
		span.SetAttributes(
			// 设置RPC方法
			semconv.RPCMethodKey.String(info.FullMethod),
			// 设置对等方的名称
			semconv.NetPeerNameKey.String(b.PeerName(ctx)),
			// 设置对等方的IP地址
			attribute.Key("net.peer.ip").String(b.PeerIP(ctx)),
		)
		// 在函数结束时设置span的状态
		defer func() {
			if err != nil {
				// 记录错误
				span.RecordError(err)
				// 如果错误类型为grpc错误，则设置span的状态码
				if e := errors.FromError(err); e != nil {
					span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(e.Code)))
				}
				// 设置span的状态为错误
				span.SetStatus(codes.Error, err.Error())
			} else {
				// 设置span的状态为成功
				span.SetStatus(codes.Ok, "OK")
			}
		}()
		// 调用handler函数
		return handler(ctx, req)
	}
}

func (b *OTELInterceptorBuilder) BuildUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	tracer := b.tracer
	if tracer == nil {
		tracer = otel.GetTracerProvider().
			Tracer("gitee.com/geekbang/basic-go/webook/pkg/grpcx")
	}
	propagator := b.propagator
	if propagator == nil {
		propagator = otel.GetTextMapPropagator()
	}
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("grpc"),
		attribute.Key("rpc.grpc.kind").String("unary"),
		attribute.Key("rpc.component").String("client"),
	}
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		newAttrs := append(attrs,
			semconv.RPCMethodKey.String(method),
			semconv.NetPeerNameKey.String(b.serviceName))
		ctx, span := tracer.Start(ctx, method,
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(newAttrs...))
		ctx = inject(ctx, propagator)
		defer func() {
			if err != nil {
				span.RecordError(err)
				if e := errors.FromError(err); e != nil {
					span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(e.Code)))
				}
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "OK")
			}
			span.End()
		}()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// 从传入的上下文中提取元数据
func extract(ctx context.Context, propagators propagation.TextMapPropagator) context.Context {
	// 从传入的上下文中获取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	// 如果没有获取到元数据，则创建一个空的元数据
	if !ok {
		md = metadata.MD{}
	}
	// 使用提取器从元数据中提取上下文
	return propagators.Extract(ctx, GrpcHeaderCarrier(md))
}

// 将元数据注入到上下文中
func inject(ctx context.Context, propagators propagation.TextMapPropagator) context.Context {
	// 从传入的上下文中获取元数据
	md, ok := metadata.FromOutgoingContext(ctx)
	// 如果没有获取到元数据，则创建一个空的元数据
	if !ok {
		md = metadata.MD{}
	}
	// 使用注入器将元数据注入到上下文中
	propagators.Inject(ctx, GrpcHeaderCarrier(md))
	// 返回新的上下文
	return metadata.NewOutgoingContext(ctx, md)
}

// GrpcHeaderCarrier ...
type GrpcHeaderCarrier metadata.MD

// Get returns the value associated with the passed key.
func (mc GrpcHeaderCarrier) Get(key string) string {
	vals := metadata.MD(mc).Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// Set stores the key-value pair.
func (mc GrpcHeaderCarrier) Set(key string, value string) {
	metadata.MD(mc).Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (mc GrpcHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range metadata.MD(mc) {
		keys = append(keys, k)
	}
	return keys
}
