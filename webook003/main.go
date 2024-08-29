package main

import (
	"github.com/gin-gonic/gin"
	"goworkwebook/webook003/ioc"
	"log"
)

func main() {
	server := gin.Default()

	db := ioc.InitDB()
	ioc.InitMiddleware(server)
	ioc.InitWeb(server, db)

	err := server.Run(":8080")
	if err != nil {
		log.Println(err)
	}
}
