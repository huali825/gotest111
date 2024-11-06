package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"goworkwebook/webook003/internal/service/sms"
	"time"
)

// Decorator 结构体，用于装饰 sms.Service 接口
type Decorator struct {
	// svc 是 sms.Service 接口的实例
	svc sms.Service
	// vector 是 prometheus.SummaryVec 类型的实例，用于记录统计信息
	vector *prometheus.SummaryVec
}

func NewDecorator(svc sms.Service, opt prometheus.SummaryOpts) *Decorator {
	return &Decorator{
		svc:    svc,
		vector: prometheus.NewSummaryVec(opt, []string{"tpl_id"}),
	}
}

// Send 函数用于发送消息
func (d *Decorator) Send(ctx context.Context,
	tplId string, args []string, numbers ...string) error {
	// 记录开始时间
	start := time.Now()
	// 在函数结束时记录耗时
	defer func() {
		duration := time.Since(start).Milliseconds()
		// 记录耗时
		d.vector.WithLabelValues(tplId).Observe(float64(duration))
	}()
	// 调用svc的Send函数发送消息
	return d.svc.Send(ctx, tplId, args, numbers...)
}
