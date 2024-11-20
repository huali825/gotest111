package opentelemetry

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	res, err := newResource("demo", "v0.0.1")
	require.NoError(t, err)

	prop := newPropagator()
	// 在客户端和服务端之间传递 tracing 的相关信息
	otel.SetTextMapPropagator(prop)

	// 初始化 trace provider //跟踪 提供者
	// 这个 provider 就是用来在打点的时候构建 trace 的
	tp, err := newTraceProvider(res)
	require.NoError(t, err)
	defer tp.Shutdown(context.Background())
	otel.SetTracerProvider(tp)

	server := gin.Default()
	// 创建一个gin的默认服务器
	server.GET("/test", func(ginCtx *gin.Context) {
		// 名字唯一
		tracer := otel.Tracer("tmh 的 openTelemetry")
		// 创建一个tracer
		var ctx context.Context = ginCtx
		// 将gin的上下文赋值给ctx
		ctx, span := tracer.Start(ctx, "top-span")
		// 开始一个span
		defer span.End()
		// 结束span
		time.Sleep(time.Second)
		// 模拟耗时操作
		span.AddEvent("发生了某事")
		// 添加一个事件
		ctx, subSpan := tracer.Start(ctx, "sub-span")
		// 开始一个子span
		defer subSpan.End()
		// 结束子span
		subSpan.SetAttributes(attribute.String("attr1", "value1"))
		// 设置子span的属性
		time.Sleep(time.Millisecond * 300)
		// 模拟耗时操作
		ginCtx.String(http.StatusOK, "测试 span")
		// 返回一个字符串
	})
	server.Run(":8082")
	// 启动服务器，监听8082端口
}

func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		))
}

func newTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	exporter, err := zipkin.New(
		"http://localhost:9411/api/v2/spans")
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
