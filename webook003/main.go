package main

import (
	"github.com/gin-gonic/gin"
	"goworkwebook/webook003/ioc"
	"log"
)

func main() {

	//server := gin.Default()
	//server.GET("/hello", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, " 2024年9月11日10:34:47  webook hello！！！")
	//})
	//_ = server.Run(":8080")

	server := gin.Default()
	db := ioc.InitDB()
	ioc.InitMiddleware(server)
	ioc.InitWeb(server, db)

	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "hello world 你来了")
	})

	err := server.Run(":8080")
	if err != nil {
		log.Println(err)
	}
}
