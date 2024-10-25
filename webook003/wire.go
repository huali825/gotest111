//go:build wireinject

package main

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

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis, ioc.InitDB, ioc.InitLogger,
		// DAO 部分
		dao.NewUserDAO,
		dao.NewArticleGORMDAO,

		// cache 部分
		cache.NewCodeCache, cache.NewUserCache,
		cache.NewArticleRedisCache,

		// repository 部分
		repository.NewUserRepository,
		repository.NewCodeRepository,
		repository.NewCachedArticleRepository,

		// Service 部分
		ioc.InitSMSService,
		ioc.InitWechatService,
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,

		// handler 部分
		web.NewUserHandler,
		ijwt.NewRedisJWTHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
