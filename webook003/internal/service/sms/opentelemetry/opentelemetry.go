package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"goworkwebook/webook003/internal/service/sms"
)

// 定义一个装饰器结构体，包含一个sms.Service和一个trace.Tracer
type Decorator struct {
	svc    sms.Service
	tracer trace.Tracer
}

// 创建一个装饰器实例，传入一个sms.Service和一个trace.Tracer
func NewDecorator(svc sms.Service, tracer trace.Tracer) *Decorator {
	return &Decorator{svc: svc, tracer: tracer}
}

// 发送短信，传入一个context.Context，一个tplId，一个args切片，一个numbers切片
func (d *Decorator) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 使用tracer.Start方法开始一个新的span
	ctx, span := d.tracer.Start(ctx, "sms")
	// 在函数结束时调用span.End方法结束span
	defer span.End()
	// 设置span的属性
	span.SetAttributes(attribute.String("tpl", tplId))
	// 添加一个事件
	span.AddEvent("发短信")
	// 调用svc.Send方法发送短信
	err := d.svc.Send(ctx, tplId, args, numbers...)
	// 如果发送短信出错，记录错误
	if err != nil {
		span.RecordError(err)
	}
	// 返回错误
	return err
}
