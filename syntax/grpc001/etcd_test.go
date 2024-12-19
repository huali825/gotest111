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
	"net"
	"testing"
	"time"
)

//=======================//
// 测试etcd 作为注册中心
//=======================//

func TestEtcd(t *testing.T) {
	suite.Run(t, new(EtcdTestSuite))
}

type EtcdTestSuite struct {
	suite.Suite
	cli *etcdv3.Client //客户端提供并管理 etcd v3客户端会话。
}

func (s *EtcdTestSuite) SetupSuite() {
	// 需要docker 开启etcd服务,并设置端口 12379
	cli, err := etcdv3.NewFromURL("localhost:12379")

	// etcdv3.NewFromURLs()
	// etcdv3.New(etcdv3.Config{Endpoints: })
	require.NoError(s.T(), err)
	s.cli = cli
}

func (s *EtcdTestSuite) TestStartServer() {

	// 监听地址 这里创建了一个server, 并监听8090端口
	listener, err := net.Listen("tcp", ":8090")
	require.NoError(s.T(), err)

	// 创建一个grpc服务器
	server := grpc.NewServer()
	// 注册UserServiceServer
	myGrpc.RegisterUserServiceServer(server, &Server{})

	// 启动服务器
	server.Serve(listener)

	// 停止服务器
	server.GracefulStop()
}

// TestETCDServer 启动 etcd 注册中心
func (s *EtcdTestSuite) TestETCDServer() {
	t := s.T()
	// 创建一个新的endpoints管理器
	em, err := endpoints.NewManager(s.cli, "service/user")
	require.NoError(t, err)

	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 定义一个地址
	addr := "127.0.0.1:8090"
	// 定义一个key
	key := "service/user/" + addr
	//// 监听地址 这里创建了一个server, 并监听8090端口
	//listener, err := net.Listen("tcp", ":8090")
	//require.NoError(s.T(), err)

	// 创建一个带有超时的上下文
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 定义租期 5秒
	var ttl int64 = 9
	// 创建一个租约
	leaseResp, err := s.cli.Grant(ctx, ttl)
	require.NoError(t, err)

	// 创建节点 service/user/127.0.0.1:8090
	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		// 定位信息，客户端怎么连你
		Addr: addr,
	}, etcdv3.WithLease(leaseResp.ID))
	require.NoError(t, err)
	//完成了注册

	// 创建一个上下文
	kaCtx, kaCancel := context.WithCancel(context.Background())
	// 启动一个协程，保持租约
	go func() {
		ch, err1 := s.cli.KeepAlive(kaCtx, leaseResp.ID)
		require.NoError(t, err1)
		for kaResp := range ch {
			t.Log(kaResp.String())
		}
	}()

	// 启动一个协程，模拟注册信息变动
	//	// 创建一个新的endpoints管理器
	//	em, err := endpoints.NewManager(s.cli, "service/user")
	//	// 定义一个地址
	//	addr := "127.0.0.1:8090"
	//	// 定义一个key
	//	key := "service/user/" + addr
	go func(endpointsManager endpoints.Manager, goFuncKey string, goFuncAddr string, goFuncLeaseResp *etcdv3.LeaseGrantResponse) {
		ticker := time.NewTicker(time.Second)
		// 循环获取ticker.C的值
		for now := range ticker.C {
			// 创建一个带有超时的context
			ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second)
			// 更新etcd中的值
			err1 := endpointsManager.Update(ctx1, []*endpoints.UpdateWithOpts{
				{
					Update: endpoints.Update{
						Op:  endpoints.Add,
						Key: goFuncKey,
						Endpoint: endpoints.Endpoint{
							Addr:     addr,
							Metadata: now.String(),
						},
					},
					Opts: []etcdv3.OpOption{etcdv3.WithLease(leaseResp.ID)},
				},
			})
			// 取消context
			cancel1()
			// 如果更新失败，打印错误信息
			if err1 != nil {
				t.Log(err1)
			}
		}
	}(em, key, addr, leaseResp)

	//// 创建一个grpc服务器
	//server := grpc.NewServer()
	//// 注册UserServiceServer
	//myGrpc.RegisterUserServiceServer(server, &Server{})
	//
	//// 启动服务器
	//server.Serve(listener)
	time.Sleep(time.Minute * 10)

	// 取消上下文
	kaCancel()
	// 删除endpoints
	err = em.DeleteEndpoint(ctx, key)
	if err != nil {
		t.Log(err)
	}

	//// 停止服务器
	//server.GracefulStop()
	// 关闭etcd客户端
	s.cli.Close()
}

// 测试Etcd客户端
func (s *EtcdTestSuite) TestClient() {
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
	// 创建带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 在函数结束时取消上下文
	defer cancel()

	// 调用GetByID方法获取用户信息
	resp, err := client.GetByID(ctx, &myGrpc.GetByIDRequest{Id: 123})
	// 断言没有错误
	require.NoError(t, err)
	// 打印用户信息
	t.Log(resp.User)
}
