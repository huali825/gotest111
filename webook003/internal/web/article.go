package web

import "github.com/gin-gonic/gin"

var _ Handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	//aSvc service.ArticleService
	//l    logger.LoggerV1
}

// l logger.LoggerV1, svc service.ArticleService
func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{}
}
func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")

	//g.PUT("/", h.Edit)
	g.POST("/edit", h.Edit)
	g.POST("/publish", h.Publish)
	g.POST("/withdraw", h.Withdraw)
}

func (h *ArticleHandler) Edit(context *gin.Context) {

}

func (h *ArticleHandler) Publish(context *gin.Context) {

}

func (h *ArticleHandler) Withdraw(context *gin.Context) {

}
