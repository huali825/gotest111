package logger

// 第一种需要用户自己使用占位符
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func example() {
	var l Logger
	l.Info("用户的微信 id %d", 123)
}

// zap风格, 日志参数都是有名字的通过一个struct来配套输出
type LoggerV1 interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
}

type Field struct {
	Key string
	Val any
}

func exampleV1() {
	var l LoggerV1
	// 这是一个新用户 union_id=123
	l.Info("这是一个新用户", Field{Key: "union_id", Val: 123})
}

// 它要去 args 必须是偶数，并且是以 key1,value1,key2,value2 的形式传递
type LoggerV2 interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}
