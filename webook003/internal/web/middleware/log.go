package middleware

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

// LogMiddlewareBuilder 是一个用于构建日志中间件的构建器
type LogMiddlewareBuilder struct {
	logFn         func(ctx context.Context, l AccessLog) // 日志函数
	allowReqBody  bool                                   // 是否允许记录请求体
	allowRespBody bool                                   // 是否允许记录响应体
}

// NewLogMiddlewareBuilder 创建一个新的 LogMiddlewareBuilder
func NewLogMiddlewareBuilder(logFn func(ctx context.Context, l AccessLog)) *LogMiddlewareBuilder {
	return &LogMiddlewareBuilder{
		logFn: logFn,
	}
}

// AllowReqBody 允许记录请求体
func (l *LogMiddlewareBuilder) AllowReqBody() *LogMiddlewareBuilder {
	l.allowReqBody = true
	return l
}

// AllowRespBody 允许记录响应体
func (l *LogMiddlewareBuilder) AllowRespBody() *LogMiddlewareBuilder {
	l.allowRespBody = true
	return l
}

// Build 构建日志中间件
func (l *LogMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if len(path) > 1024 {
			path = path[:1024]
		}
		method := ctx.Request.Method
		al := AccessLog{
			Path:   path,
			Method: method,
		}
		if l.allowReqBody {
			// Request.Body 是一个 Stream 对象，只能读一次
			body, _ := ctx.GetRawData()
			if len(body) > 2048 {
				al.ReqBody = string(body[:2048])
			} else {
				al.ReqBody = string(body)
			}
			// 放回去
			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
			//ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		start := time.Now()

		if l.allowRespBody {
			ctx.Writer = &responseWriter{
				ResponseWriter: ctx.Writer,
				al:             &al,
			}
		}

		defer func() {
			al.Duration = time.Since(start)
			//duration := time.Now().Sub(start)
			l.logFn(ctx, al)
		}()

		// 直接执行下一个 middleware...直到业务逻辑
		ctx.Next()
		// 在这里，你就拿到了响应
	}
}

// AccessLog 访问日志结构体
type AccessLog struct {
	Path     string        `json:"path"`      // 请求路径
	Method   string        `json:"method"`    // 请求方法
	ReqBody  string        `json:"req_body"`  // 请求体
	Status   int           `json:"status"`    // 响应状态码
	RespBody string        `json:"resp_body"` // 响应体
	Duration time.Duration `json:"duration"`  // 请求耗时
}

// responseWriter 自定义的响应写入器
type responseWriter struct {
	gin.ResponseWriter
	al *AccessLog
}

// Write 写入响应体
func (w *responseWriter) Write(data []byte) (int, error) {
	w.al.RespBody = string(data)
	return w.ResponseWriter.Write(data)
}

// WriteHeader 写入响应状态码
func (w *responseWriter) WriteHeader(statusCode int) {
	w.al.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
