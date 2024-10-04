package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"goworkwebook/webook003/internal/repository"
	"goworkwebook/webook003/internal/repository/cache"
	"goworkwebook/webook003/internal/repository/dao"
	"goworkwebook/webook003/internal/service"
	"goworkwebook/webook003/internal/web"
	"goworkwebook/webook003/ioc"
)

func InitWeb(server *gin.Engine, db *gorm.DB) *web.UserHandler {
	cmdable := ioc.InitRedis()
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(cmdable)
	codeCache := cache.NewCodeCache(cmdable)

	repo := repository.NewUserRepository(userDAO, userCache)
	codeRep := repository.NewCodeRepository(codeCache)

	smsService := ioc.InitSMSService()
	svc := service.NewUserService(repo)
	codeS := service.NewCodeService(codeRep, smsService)

	u := web.NewUserHandler(svc, codeS)
	u.RegisterRoutes(server)
	return u
}
