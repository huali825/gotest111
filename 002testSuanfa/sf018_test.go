/*
date:2025年4月24日10:57:32
title: 2.2作业
author:tmh
*/

package _02testSuanfa

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func Test018Name(t *testing.T) {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin路由器
	router := gin.New()

	// 使用自定义的日志中间件
	router.Use(loggerMiddleware())

	// 设置路由
	router.Any("/*path", func(c *gin.Context) {
		// 1. 将请求头写入响应头
		for key, values := range c.Request.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		// 2. 读取环境变量VERSION并写入响应头
		version := os.Getenv("VERSION")
		if version != "" {
			c.Header("VERSION", version)
		}

		// 返回200状态码
		c.Status(http.StatusOK)
	})

	// 健康检查路由
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// 启动服务器
	fmt.Println("Server is running on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// 自定义日志中间件
func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 计算请求耗时
		latency := endTime.Sub(startTime)

		// 获取客户端IP
		clientIP := c.ClientIP()

		// 获取HTTP状态码
		statusCode := c.Writer.Status()

		// 记录日志到标准输出
		log.Printf("[%s] %s %s %d %v",
			clientIP,
			c.Request.Method,
			c.Request.URL.Path,
			statusCode,
			latency,
		)
	}
}

func Test018to01Name(t *testing.T) {

}
