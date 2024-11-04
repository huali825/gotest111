package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	//server := gin.Default()
	//server.GET("/hello", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, " 2024年9月11日10:34:47  webook hello！！！")
	//})
	//_ = server.Run(":8080")

	initViperV1() // 初始化配置
	initLogger()  // 初始化日志

	app := InitWebServerAndCsm()
	initPrometheus()

	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	server := app.server
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello，启动成功了！")
	})

	err := server.Run(":8080")
	if err != nil {
		log.Println(err)
	}
}

// 初始化 Prometheus
func initPrometheus() {
	go func() {
		// 专门给 prometheus 用的端口
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}

// 有用 但是我暂时没用到这里
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

// 初始化日志 目前暂时没用了2024年11月4日16:29:00
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
