package prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type Builder struct {
	Namespace  string
	Subsystem  string
	Name       string
	InstanceId string
	Help       string
}

// Namespace：指标命名空间，通常用于标识项目或团队。
// Subsystem：子系统名称，用于进一步细分指标。
// Name：指标名称。
// InstanceId：实例ID，用于区分不同的服务实例。
// Help：指标的描述信息。

func (b *Builder) BuildResponseTime() gin.HandlerFunc {
	// 定义标签
	labels := []string{"method", "pattern", "status"}

	// 创建一个新的Summary指标
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Help:      b.Help,
		Name:      b.Name + "_resp_time",
		ConstLabels: map[string]string{
			"instance_id": b.InstanceId,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, labels)

	// 注册指标
	prometheus.MustRegister(vector)

	// 返回Gin中间件函数
	return func(ctx *gin.Context) {
		// 记录开始时间
		start := time.Now()

		// 在请求处理完成后执行
		defer func() {
			// 计算响应时间
			duration := time.Since(start).Milliseconds()
			// 获取请求方法和路由
			method := ctx.Request.Method
			pattern := ctx.FullPath()
			// 获取HTTP状态码
			status := ctx.Writer.Status()
			// 上报指标
			vector.WithLabelValues(
				method, pattern, strconv.Itoa(status)).
				Observe(float64(duration))
		}()

		// 继续处理请求
		ctx.Next()
	}
}

func (b *Builder) BuildActiveRequest() gin.HandlerFunc {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Help:      b.Help,
		// Namespace 和 Subsystem 和 Name 都不能有 _ 以外的其它符号
		Name: b.Name + "_active_req",
		ConstLabels: map[string]string{
			"instance_id": b.InstanceId,
		},
	})
	prometheus.MustRegister(gauge)
	return func(ctx *gin.Context) {
		gauge.Inc()
		defer gauge.Dec()
		ctx.Next()
	}
}
