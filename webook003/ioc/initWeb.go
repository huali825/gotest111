package ioc

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"goworkwebook/webook003/internal/repository"
	"goworkwebook/webook003/internal/repository/dao"
	"goworkwebook/webook003/internal/service"
	"goworkwebook/webook003/internal/web"
)

func InitWeb(server *gin.Engine, db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(server, svc)
	u.RegisterRoutes()
	return u
}
