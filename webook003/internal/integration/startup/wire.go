//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"goworkwebook/webook003/internal/repository"
	"goworkwebook/webook003/internal/repository/cache"
	"goworkwebook/webook003/internal/repository/dao"
	"goworkwebook/webook003/internal/service"
	"goworkwebook/webook003/internal/web"
	ijwt "goworkwebook/webook003/internal/web/jwt"
	"goworkwebook/webook003/ioc"
)

// thirdPartySet 定义了第三方依赖
var thirdPartySet = wire.NewSet( // 第三方依赖
	InitRedis, InitDB2,
	InitLogger)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		//InitRedis, InitDB2, InitLogger,
		thirdPartySet,

		// DAO 部分
		dao.NewUserDAO,

		// cache 部分
		cache.NewCodeCache, cache.NewUserCache,

		// repository 部分
		repository.NewUserRepository,
		repository.NewCodeRepository,

		// Service 部分
		ioc.InitSMSService,
		ioc.InitWechatService,
		service.NewUserService,
		service.NewCodeService,

		// handler 部分
		web.NewUserHandler,
		ijwt.NewRedisJWTHandler,
		web.NewOAuth2WechatHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}

// InitArticleHandler 初始化ArticleHandler
func InitArticleHandler() *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		dao.NewArticleGORMDAO,
		service.NewArticleService,
		web.NewArticleHandler,
		repository.NewCachedArticleRepository)
	return &web.ArticleHandler{}
}
