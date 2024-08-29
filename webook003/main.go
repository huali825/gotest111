package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"goworkwebook/webook003/ioc"
	"log"
	"strings"
	"time"
)

func main() {
	server := gin.Default()
	db := ioc.InitDB()
	server.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		//ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "localhost")
		},
	}))

	//server.GET("/hello", func(c *gin.Context) {
	//	c.String(http.StatusOK, "Hello, World!")
	//})
	//server.GET("/ping", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"message": "pong",
	//	})
	//})

	ioc.InitWeb(server, db)

	err := server.Run(":8080")
	if err != nil {
		log.Println(err)
	}
}
