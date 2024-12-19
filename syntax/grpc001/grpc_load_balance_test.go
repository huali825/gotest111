package grpc001

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	myGrpc "goworkwebook/syntax/grpc001/myGrpc"
	"goworkwebook/webook003/pkg/netx"
	"net"
	"testing"
	"time"
)

func TestBalance001(t *testing.T) {
	suite.Run(t, new(BalancerTestSuite))
}

func (s *BalancerTestSuite) SetupSuite() {
	cli, err := etcdv3.NewFromURL("localhost:12379")
	// etcdv3.NewFromURLs()
	// etcdv3.New(etcdv3.Config{Endpoints: })
	require.NoError(s.T(), err)
	s.cli = cli
}

type BalancerTestSuite struct {
	suite.Suite
	cli *etcdv3.Client
}

func (s *BalancerTestSuite) TestPickFirst() {
	// 获取测试用例的T对象
	t := s.T()
	// 创建Etcd解析器
	etcdResolver, err := resolver.NewBuilder(s.cli)
	// 断言没有错误
	require.NoError(s.T(), err)

	// 使用Etcd解析器创建gRPC连接
	cc, err := grpc.Dial("etcd:///service/user",
		grpc.WithResolvers(etcdResolver),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	// 断言没有错误
	require.NoError(t, err)

	// 创建UserService客户端
	client := myGrpc.NewUserServiceClient(cc)

	// 调用GetByID方法获取用户信息
	for i := 0; i < 10; i++ {
		// 创建带有超时的上下文
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		// 在函数结束时取消上下文
		defer cancel()
		resp, err := client.GetByID(ctx, &myGrpc.GetByIDRequest{Id: 123})
		// 断言没有错误
		require.NoError(t, err)
		// 打印用户信息
		t.Log(resp.User)
	}

	//time.Sleep(time.Minute)
}

func (s *BalancerTestSuite) TestServer() {
	go func() {
		s.startServer(":8090")
	}()
	s.startServer(":8091")
}

func (s *BalancerTestSuite) startServer(addr string) {
	t := s.T()

	// 监听指定地址
	listenS, err := net.Listen("tcp", addr)
	require.NoError(s.T(), err)

	// 创建etcd管理器
	em, err := endpoints.NewManager(s.cli, "service/user")
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//addr := "127.0.0.1:8090"
	// 获取本机ip 互联网视角
	addr = netx.GetOutboundIP() + addr
	// 构建key
	key := "service/user/" + addr
	t.Log(addr, key)

	// 创建租期
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var ttl int64 = 5
	leaseResp, err := s.cli.Grant(ctx, ttl)
	require.NoError(t, err)

	// 添加endpoint
	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		// 定位信息，客户端怎么连你
		Addr: addr,
	}, etcdv3.WithLease(leaseResp.ID))
	require.NoError(t, err)
	// 创建keepalive上下文
	kaCtx, kaCancel := context.WithCancel(context.Background())
	// 启动keepalive协程
	go func() {
		_, err1 := s.cli.KeepAlive(kaCtx, leaseResp.ID)
		require.NoError(t, err1)
		//for kaResp := range ch {
		//	t.Log(kaResp.String())
		//}
	}()

	// 启动模拟注册信息变动的协程

	// 创建grpc服务器
	server := grpc.NewServer()
	// 注册UserServiceServer
	myGrpc.RegisterUserServiceServer(server, &Server{Name: addr})
	// 启动grpc服务器
	server.Serve(listenS)
	// 取消keepalive上下文
	kaCancel()
	// 删除endpoint
	err = em.DeleteEndpoint(ctx, key)
	if err != nil {
		t.Log(err)
	}
	// 停止grpc服务器
	server.GracefulStop()
	// 关闭etcd客户端
	s.cli.Close()
}