package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"goworkwebook/webook003/internal/web/middleware"
	"strings"
	"time"
)

func InitMiddleware(server *gin.Engine) *gin.Engine {
	server.Use(func(context *gin.Context) {
		println("这是我的 Middleware 1")
	})

	server.Use(
		func(ctx *gin.Context) {
			println("这是我的 Middleware 2")
		},
		cors.New(cors.Config{
			//AllowOrigins:     []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
			//ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					return true
				}
				return strings.Contains(origin, "localhost")
			},
		}),
		func(ctx *gin.Context) {
			println("这是我的 Middleware 3")
		},
	)

	server.Use(func(context *gin.Context) {
		println("这是我的 Middleware 4")
	})

	//session搞好了
	// 创建一个cookie存储
	store := cookie.NewStore([]byte("secret"))
	// 使用sessions中间件，将cookie存储命名为"mysession"
	server.Use(
		sessions.Sessions("mysession", store),
	)
	server.Use(
		middleware.NewLoginMiddlewareBuilder().
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").Build(),
	)

	return server
}
