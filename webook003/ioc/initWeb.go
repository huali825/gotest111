package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"goworkwebook/webook003/internal/web"
	ijwt "goworkwebook/webook003/internal/web/jwt"
	"goworkwebook/webook003/internal/web/middleware"
	"goworkwebook/webook003/pkg/ginx/middleware/prometheus"
	"goworkwebook/webook003/pkg/logger"
	"strings"
	"time"
)

func InitWebServer(
	mdls []gin.HandlerFunc,
	userHdl *web.UserHandler,
	artHdl *web.ArticleHandler,
	wechatHdl *web.OAuth2WechatHandler) *gin.Engine {

	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	artHdl.RegisterRoutes(server)
	wechatHdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable, hdl ijwt.Handler, l logger.LoggerV1) []gin.HandlerFunc {
	// func
	pb := &prometheus.Builder{
		Namespace: "tangminghao_zixue", //不要用连字符 -
		Subsystem: "webook",
		Name:      "gin_http",
		Help:      "统计 GIN 的HTTP接口数据",
	}
	//ginx.InitCounter(prometheus2.CounterOpts{
	//	Namespace: "geektime_daming",
	//	Subsystem: "webook",
	//	Name:      "biz_code",
	//	Help:      "统计业务错误码",
	//})
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
			ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
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

		//func(context *gin.Context) {
		//	println("redis限流, pkg/ginx/middleware/ratelimit实现的 集成测试需要关闭此功能")
		//},
		//ratelimit.NewBuilder(redisClient, time.Second, 1).Build(),

		func(context *gin.Context) { println("jwt登录校验") },
		middleware.NewLoginJWTMiddlewareBuilder(hdl).
			IgnorePaths("/hello").
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/oauth2/wechat/authurl").
			IgnorePaths("/oauth2/wechat/callback").
			IgnorePaths("/users/login_sms").CheckLogin(),

		//在处理HTTP请求和响应时记录日志
		//middleware.NewLogMiddlewareBuilder(func(ctx context.Context, al middleware.AccessLog) {
		//	l.Debug("", logger.Field{Key: "req", Val: al})
		//}).AllowReqBody().AllowRespBody().Build(),

		pb.BuildResponseTime(),

		//gin框架集成open telemetry
		otelgin.Middleware("webook"),

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
}
