package grpcInterceptor

/*
 * @Date: 2024年12月27日
 * @LastEditors: TMH
 * @LastEditTime: 2024年12月29日12:51:27
 * @FilePath: syntax/002grpcInterceptor/client_test.go
 * @Description: grpc 客户端 多个拦截器(类似middleware)的test实现
 */

import (
	"context"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
)

// 测试客户端
func TestClient(t *testing.T) {
	// 连接服务器
	cc, err := grpc.NewClient("localhost:8090",
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	//cc, err := grpc.Dial("localhost:8090", grpc.WithInsecure())
	require.NoError(t, err)
	// 创建客户端
	client := NewUserServiceClient(cc)
	// 调用客户端方法
	resp, err := client.GetByID(context.Background(), &GetByIDRequest{Id: 123})
	require.NoError(t, err)
	// 打印返回结果
	t.Log(resp.User)
}

// 测试客户端 限流
func TestClient02(t *testing.T) {
	// 连接服务器
	//cc, err := grpc.Dial("localhost:8090", grpc.WithInsecure())
	cc, err := grpc.NewClient("localhost:8090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(clientFirst, clientSecond))
	require.NoError(t, err)

	// 创建客户端
	client := NewUserServiceClient(cc)
	// 调用客户端方法
	resp, err := client.GetByID(context.Background(), &GetByIDRequest{Id: 123})
	require.NoError(t, err)
	// 打印返回结果
	t.Log(resp.User)
}

// clientFirst 限流拦截器1
var clientFirst grpc.UnaryClientInterceptor = func(ctx context.Context,
	method string, req, reply any,
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	log.Println("before call")
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Println("after call")
	return err
}

// clientSecond 限流拦截器2
var clientSecond grpc.UnaryClientInterceptor = func(ctx context.Context,
	method string, req, reply any,
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	log.Println("before call2")
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Println("after call2")
	return err
}
