package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	grpc2 "goworkwebook/004interactive/grpc"
	"goworkwebook/webook003/pkg/grpcx"
	"goworkwebook/webook003/pkg/logger"
)

// NewGrpcxServer 创建一个新的grpcx.Server实例
func NewGrpcxServer(intrSvc *grpc2.InteractiveServiceServer, l logger.LoggerV1) *grpcx.Server {
	// 定义Config结构体，用于存储配置信息
	type Config struct {
		EtcdAddr string `yaml:"etcdAddr"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	}
	// 创建一个新的grpc.Server实例
	s := grpc.NewServer()
	// 将InteractiveServiceServer注册到grpc.Server实例中
	intrSvc.Register(s)
	// 定义一个Config实例
	var cfg Config
	// 从viper中获取grpc.server的配置信息，并存储到cfg中
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		// 如果获取配置信息失败，则panic
		panic(err)
	}
	// 返回一个新的grpcx.Server实例
	return &grpcx.Server{
		Server:   s,
		EtcdAddr: cfg.EtcdAddr,
		Port:     cfg.Port,
		Name:     cfg.Name,
		L:        l,
	}
}
