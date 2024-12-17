package gozero

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"testing"
	"time"
)

type GoZeroTestSuite struct {
	suite.Suite
}

func TestGoZero(t *testing.T) {
	suite.Run(t, new(GoZeroTestSuite))
}

// TestGoZeroClient 启动 grpc 客户端
func (s *GoZeroTestSuite) TestGoZeroClient() {
	// 创建 grpc 客户端
	zClient := zrpc.MustNewClient(
		zrpc.RpcClientConf{
			Etcd: discov.EtcdConf{
				Hosts: []string{"localhost:12379"},
				Key:   "user",
			},
		},
		zrpc.WithDialOption(
			grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		))

	// 创建 UserServiceClient
	client := NewUserServiceClient(zClient.Conn())

	// 循环调用 GetByID 方法
	for i := 0; i < 10; i++ {
		// 创建上下文
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		// 调用 GetByID 方法
		resp, err := client.GetByID(ctx, &GetByIDRequest{
			Id: 123,
		})
		// 取消上下文
		cancel()
		// 断言 err 为 nil
		require.NoError(s.T(), err)
		// 打印 resp.User
		s.T().Log(resp.User)
	}
}

// TestGoZeroServer 启动 grpc 服务端
func (s *GoZeroTestSuite) TestGoZeroServer() {
	go func() {
		s.startServer(":8090")
	}()
	s.startServer(":8091")
}

// 启动服务器
func (s *GoZeroTestSuite) startServer(addr string) {

	// 创建RPC服务器配置
	c := zrpc.RpcServerConf{
		ListenOn: addr, // 监听地址
		Etcd: discov.EtcdConf{
			Hosts: []string{"localhost:12379"}, // Etcd主机地址
			Key:   "user",                      // Etcd键
		},
	}

	// 创建RPC服务器
	server := zrpc.MustNewServer(c, func(grpcServer *grpc.Server) {
		// 注册UserServiceServer
		RegisterUserServiceServer(grpcServer, &Server{
			Name: addr,
		})
	})

	// 启动服务器
	server.Start()
}
