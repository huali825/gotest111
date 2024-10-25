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

var userSvcProvider = wire.NewSet(
	dao.NewUserDAO,
	cache.NewUserCache,
	repository.NewCachedUserRepository,
	service.NewUserService)

var articleSvcProvider = wire.NewSet(
	dao.NewArticleGORMDAO,
	cache.NewArticleRedisCache,
	repository.NewCachedArticleRepository,
	service.NewArticleService)

var interactiveSvcSet = wire.NewSet(
	dao.NewGORMInteractiveDAO,
	cache.NewInteractiveRedisCache,
	repository.NewCachedInteractiveRepository,
	service.NewInteractiveService,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		//InitRedis, InitDB2, InitLogger,
		thirdPartySet,
		userSvcProvider,
		articleSvcProvider,
		interactiveSvcSet,

		// cache 部分
		cache.NewCodeCache,

		// repository 部分
		repository.NewCodeRepository,

		// Service 部分
		ioc.InitSMSService,
		InitWechatService,
		//ioc.InitWechatService,
		service.NewCodeService,

		// handler 部分
		web.NewUserHandler,
		web.NewArticleHandler,
		//web.interactiveHandler, //在 NewArticleHandler 中初始化

		ijwt.NewRedisJWTHandler,
		web.NewOAuth2WechatHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}

// InitArticleHandler 初始化ArticleHandler
func InitArticleHandler(dao dao.ArticleDAO) *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		userSvcProvider,
		interactiveSvcSet,
		repository.NewCachedArticleRepository,
		cache.NewArticleRedisCache,
		service.NewArticleService,
		web.NewArticleHandler)
	return &web.ArticleHandler{}
}

func InitInteractiveService() service.InteractiveService {
	wire.Build(thirdPartySet, interactiveSvcSet)
	return service.NewInteractiveService(nil)
}
