//go:build wireinject

package main

import (
	"github.com/google/wire"
	"goworkwebook/webook003/internal/events/article"
	"goworkwebook/webook003/internal/repository"
	"goworkwebook/webook003/internal/repository/cache"
	"goworkwebook/webook003/internal/repository/dao"
	"goworkwebook/webook003/internal/service"
	"goworkwebook/webook003/internal/web"
	ijwt "goworkwebook/webook003/internal/web/jwt"
	"goworkwebook/webook003/ioc"
)

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

var rankingSvcSet = wire.NewSet(
	cache.NewRankingRedisCache,
	repository.NewCachedRankingRepository,
	service.NewBatchRankingService,
)

func InitWebServerAndCsm() *App {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis, ioc.InitDB, ioc.InitLogger,

		// Sarama
		ioc.InitSaramaClient,
		ioc.InitSyncProducer,
		ioc.InitConsumers,
		article.NewSaramaSyncProducer,
		article.NewInteractiveReadEventConsumer,

		userSvcProvider,
		articleSvcProvider,
		interactiveSvcSet,
		rankingSvcSet,
		ioc.InitJobs,
		ioc.InitRankingJob,

		// cache 部分
		cache.NewCodeCache,

		// repository 部分
		repository.NewCodeRepository,

		// Service 部分
		ioc.InitSMSService,
		ioc.InitWechatService,
		//ioc.InitWechatService,
		service.NewCodeService,

		// handler 部分
		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WechatHandler,
		ijwt.NewRedisJWTHandler,
		//web.interactiveHandler, //在 NewArticleHandler 中初始化

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,

		// 使用wire包将App结构体注入到wireApp变量中
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
