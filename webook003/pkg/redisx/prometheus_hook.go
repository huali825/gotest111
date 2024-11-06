package redisx

import (
	"context"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"net"
	"strconv"
	"time"
)

func use(client *redis.Client) {
	client.AddHook(NewPrometheusHook(prometheus.SummaryOpts{
		Name: "redis_process_duration",
		Help: "redis process duration",
		// 定义SummaryVec的标签
		ConstLabels: map[string]string{
			"addr": client.Options().Addr,
		},
	}))
}

//type Hook interface {
//    DialHook(next DialHook) DialHook
//    ProcessHook(next ProcessHook) ProcessHook
//    ProcessPipelineHook(next ProcessPipelineHook) ProcessPipelineHook
//}

// PrometheusHook 定义一个PrometheusHook结构体，包含一个指向prometheus.SummaryVec的指针
type PrometheusHook struct {
	vector *prometheus.SummaryVec
}

// NewPrometheusHook 创建一个新的PrometheusHook实例，并返回指向该实例的指针
func NewPrometheusHook(opt prometheus.SummaryOpts) *PrometheusHook {
	// 使用传入的SummaryOpts创建一个新的SummaryVec，并将其赋值给PrometheusHook的vector字段
	return &PrometheusHook{
		vector: prometheus.NewSummaryVec(opt, []string{"cmd", "key_exist"}),
	}
}

// DialHook 函数用于创建一个新的redis.DialHook，该函数接收一个redis.DialHook作为参数，并返回一个新的redis.DialHook
func (p *PrometheusHook) DialHook(next redis.DialHook) redis.DialHook {
	// 返回一个新的redis.DialHook，该函数接收一个context.Context、一个网络类型和一个地址作为参数，并返回一个net.Conn和一个error
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		// 调用next函数，并返回其结果
		return next(ctx, network, addr)
	}
}

// ProcessHook 函数用于创建一个新的redis.ProcessHook，该函数接收一个redis.ProcessHook作为参数，并返回一个新的redis.ProcessHook
func (p *PrometheusHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	// 返回一个新的redis.ProcessHook，该函数接收一个context.Context和一个redis.Cmder作为参数，并返回一个error
	return func(ctx context.Context, cmd redis.Cmder) error {
		// 记录开始时间
		start := time.Now()
		var err error
		// 使用defer关键字，在函数返回前执行
		defer func() {
			// 计算执行时间
			duration := time.Since(start).Milliseconds()
			// 判断错误是否为redis.Nil
			keyExists := errors.Is(redis.Nil, err)
			// 使用Prometheus记录执行时间
			p.vector.WithLabelValues(cmd.Name(), strconv.FormatBool(keyExists)).
				Observe(float64(duration))
		}()
		// 调用next函数，并返回其结果
		err = next(ctx, cmd)
		return err
	}
}

// ProcessPipelineHook 函数用于创建一个新的redis.ProcessPipelineHook，该函数接收一个redis.ProcessPipelineHook作为参数，并返回一个新的redis.ProcessPipelineHook
func (p *PrometheusHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	// 返回一个新的redis.ProcessPipelineHook，该函数接收一个context.Context和一个redis.Cmder的切片作为参数，并返回一个error
	return func(ctx context.Context, cmds []redis.Cmder) error {
		// 调用next函数，并返回其结果
		return next(ctx, cmds)
	}
}
