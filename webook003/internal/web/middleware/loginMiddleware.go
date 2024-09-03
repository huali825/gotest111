package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// 在这里验证 登录后的userid的值
		// 获取当前会话
		sess := sessions.Default(ctx)
		// 如果会话中没有userId，则返回未授权状态
		if sess.Get("userId") == nil {
			// 终止请求并返回未授权状态
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
