package web

import (
	"github.com/gin-gonic/gin"
	"goworkwebook/webook003/internal/domain"
	"goworkwebook/webook003/internal/service"
	"goworkwebook/webook003/internal/web/jwt"
	"goworkwebook/webook003/pkg/logger"
	"net/http"
)

var _ Handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc service.ArticleService
	l   logger.LoggerV1
}

//  l logger.LoggerV1, svc service.ArticleService

func NewArticleHandler(l logger.LoggerV1, svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		l:   l,
		svc: svc,
	}
}
func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")

	//g.PUT("/", h.Edit)
	g.POST("/edit", h.Edit)
	g.POST("/publish", h.Publish)
	g.POST("/withdraw", h.Withdraw)
}

// Edit 函数用于编辑文章
func (h *ArticleHandler) Edit(ctx *gin.Context) {
	// 定义请求结构体
	type Req struct {
		Id      int64
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	// 声明请求结构体变量
	var req Req
	// 绑定请求参数到结构体变量
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 从上下文中获取用户信息
	uc := ctx.MustGet("user").(jwt.UserClaims)
	// 调用服务层保存文章数据
	id, err := h.svc.Save(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	// 如果保存失败，返回错误信息
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		//使用Error的示例
		h.l.Error("保存文章数据失败",
			logger.Field{Key: "uid", Val: uc.Uid}, logger.Field{Key: "err", Val: err})
		h.l.Error("保存文章数据失败",
			logger.Int64("uid", uc.Uid), logger.Error(err))
		return
	}
	// 返回保存成功的信息
	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})

}

func (h *ArticleHandler) Publish(context *gin.Context) {

}

func (h *ArticleHandler) Withdraw(context *gin.Context) {

}
