package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
)

func main() {

	//server := gin.Default()
	//server.GET("/hello", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, " 2024年9月11日10:34:47  webook hello！！！")
	//})
	//_ = server.Run(":8080")

	initViperV1() // 初始化配置
	initLogger()  // 初始化日志

	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "hello world 你来了")
		zap.L().Info("hello world")
	})
	err := server.Run(":8080")
	if err != nil {
		log.Println(err)
	}
}

func initViperV1() {
	cfile := pflag.String("config",
		"webook003/config/config.yaml", "配置文件路径")
	// 这一步之后，cfile 里面才有值
	pflag.Parse()
	//viper.Set("db.dsn", "localhost:3306")
	// 所有的默认值放好s
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfile)
	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	val := viper.Get("test.key")
	log.Println(val)
}

// 初始化日志
func initLogger() {
	// 创建一个新的开发模式的日志
	logger, err := zap.NewDevelopment()
	// 如果创建失败，则抛出异常
	if err != nil {
		panic(err)
	}
	// 将全局日志替换为新的日志
	zap.ReplaceGlobals(logger)
}
