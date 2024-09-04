package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
	// 用 Go 的方式编码解码
	gob.Register(time.Now())
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
		id := sess.Get("userId")
		if id == nil {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		updateTime := sess.Get("update_Time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 20, //单位秒
		})

		now := time.Now()
		if updateTime == nil {
			sess.Set("update_Time", now)
			if err := sess.Save(); err != nil {
				panic(err)
			}
			println("设置update_time")
			return
		}
		//1000*60*60*24 // 24小时
		updateTimeVal, _ := updateTime.(time.Time)
		if now.Sub(updateTimeVal) > time.Second*10 {
			sess.Set("update_Time", now)
			println("刷新update_Time")
			if err := sess.Save(); err != nil {
				panic(err)
			}
		}
	}
}
