package grpcInterceptor

/*
 * @Date: 2024年12月27日
 * @LastEditors: TMH
 * @LastEditTime: 2024年12月29日12:51:27
 * @FilePath: syntax/002grpcInterceptor/server_test.go
 * @Description: grpc 服务端 多个拦截器(类似middleware)的logger test实现
 */

import (
	"context"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	loggerI "goworkwebook/webook003/pkg/grpcx/interceptors/logger"
	"goworkwebook/webook003/pkg/logger"
	"log"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(first, second, strIB.BuildServerUnaryInterceptor())) // 添加拦截器 可添加多个拦截器组成拦截链
	defer func() {
		server.GracefulStop()
	}()
	userServer := &Server{}
	RegisterUserServiceServer(server, userServer)

	//l, err := net.Listen("tcp", ":8090")
	l, err := net.Listen("tcp", "127.0.0.1:8090")
	if err != nil {
		panic(err)
	}
	t.Log(l.Addr().String())

	err = server.Serve(l)
	t.Log(err)
}

var first grpc.UnaryServerInterceptor = func(ctx context.Context,
	req any, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {
	log.Println("before call")
	resp, err = handler(ctx, req)
	log.Println("after call")
	return
}

var second grpc.UnaryServerInterceptor = func(ctx context.Context,
	req any, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {
	log.Println("before call2")
	resp, err = handler(ctx, req)
	log.Println("after call2")
	return
}

var strIB = loggerI.NewInterceptorBuilder(InitLogger())

func InitLogger() logger.LoggerV1 {
	cfg := zap.NewDevelopmentConfig()
	err := viper.UnmarshalKey("log", &cfg)
	if err != nil {
		panic(err)
	}
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}
