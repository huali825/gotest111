package web

import "github.com/gin-gonic/gin"

type Handler interface {
	RegisterRoutes(server *gin.Engine)
}

// Page 结构体用于分页查询数据
// 这段代码定义了一个名为Page的结构体，其中包含两个字段：Limit和Offset。
// Limit表示每页显示的最大数量，Offset表示当前页面的起始位置。
// 这个结构体通常用于分页查询数据时使用。
type Page struct {
	Limit  int
	Offset int
}
