package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"goworkwebook/webook003/internal/web"
	"goworkwebook/webook003/internal/web/middleware"
	"strings"
	"time"
)

func InitWebServer(
	mdls []gin.HandlerFunc,
	userHdl *web.UserHandler,
	wechatHdl *web.OAuth2WechatHandler) *gin.Engine {

	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	wechatHdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		func(context *gin.Context) {
			println("这是我的 Middleware 1")
		},

		cors.New(cors.Config{
			// 允许的来源
			//AllowOrigins: []string{"*"},
			// 允许的方法
			AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			// 允许的头部
			AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
			// 暴露的头部
			ExposeHeaders: []string{"x-jwt-token"},
			// 是否允许发送凭证
			AllowCredentials: true,
			// 最大缓存时间
			MaxAge: 12 * time.Hour,
			// 允许的来源函数
			AllowOriginFunc: func(origin string) bool {
				// 如果来源以http://localhost开头，则允许
				if strings.HasPrefix(origin, "http://localhost") {
					return true
				}
				// 如果来源包含localhost，则允许
				return strings.Contains(origin, "localhost")
			},
		}),

		//func(context *gin.Context) { println("redis限流, pkg/ginx/middleware/ratelimit实现的") },
		//ratelimit.NewBuilder(redisClient, time.Second, 1).Build(),

		func(context *gin.Context) { println("jwt登录校验") },
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").CheckLogin(),
	}

	//store的三种实现方式:
	// 第1种实现方式
	//store := cookie.NewStore([]byte("secret"))

	//第2种实现方式
	//store := memstore.NewStore([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), []byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"))

	//第3种实现方式
	//store, err := redis.NewStore(16,
	//"tcp", "localhost:6379", "",
	//[]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), []byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"))
	//if err != nil {
	//	panic(err)
	//}

	// 使用sessions中间件，将cookie存储命名为"mysession"
	//server.Use(sessions.Sessions("mysession", store))
	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/users/signup").
	//	IgnorePaths("/users/login").Build())

}
