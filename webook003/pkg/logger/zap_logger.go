package logger

import "go.uber.org/zap"

type ZapLogger struct {
	l *zap.Logger
}

// 下面的结构体方法的实现是 实现的 LoggerV1 接口
// zap风格, 日志参数都是有名字的通过一个struct来配套输出
//type LoggerV1 interface {
//	Debug(msg string, args ...Field)
//	Info(msg string, args ...Field)
//	Warn(msg string, args ...Field)
//	Error(msg string, args ...Field)
//}
//
//type Field struct {
//	Key string
//	Val any
//}
//
//func exampleV1() {
//	var l LoggerV1
//	// 这是一个新用户 union_id=123
//	l.Info("这是一个新用户", Field{Key: "union_id", Val: 123})
//}

func NewZapLogger(l *zap.Logger) *ZapLogger {
	return &ZapLogger{
		l: l,
	}
}

func (z *ZapLogger) Debug(msg string, args ...Field) {
	z.l.Debug(msg, z.toArgs(args)...)
}

func (z *ZapLogger) Info(msg string, args ...Field) {
	z.l.Info(msg, z.toArgs(args)...)
}

func (z *ZapLogger) Warn(msg string, args ...Field) {
	z.l.Warn(msg, z.toArgs(args)...)
}

func (z *ZapLogger) Error(msg string, args ...Field) {
	z.l.Error(msg, z.toArgs(args)...)
}

// 将ZapLogger的args转换为zap.Field
func (z *ZapLogger) toArgs(args []Field) []zap.Field {
	// 创建一个zap.Field切片，长度为args的长度
	res := make([]zap.Field, 0, len(args))
	// 遍历args
	for _, arg := range args {
		// 将arg的Key和Val转换为zap.Any类型，并添加到res中
		res = append(res, zap.Any(arg.Key, arg.Val))
	}
	// 返回res
	return res
}
