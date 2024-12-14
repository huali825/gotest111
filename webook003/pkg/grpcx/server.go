package grpcx

import (
	"context"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
	"goworkwebook/webook003/pkg/logger"
	"goworkwebook/webook003/pkg/netx"
	"net"
	"strconv"
	"time"
)

// Server 结构体，包含 grpc.Server、EtcdAddr、Port、Name、L、client、kaCancel 字段
type Server struct {
	*grpc.Server
	EtcdAddr string
	Port     int
	Name     string
	L        logger.LoggerV1

	client   *etcdv3.Client
	kaCancel func()
}

//func NewServer(c *etcdv3.Client) *Server {
//	return &Server{
//		client: c,
//	}
//}

// Serve 方法，用于启动服务
func (s *Server) Serve() error {
	// 获取监听地址
	addr := ":" + strconv.Itoa(s.Port)
	// 监听地址
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// 我们要在这里完成注册
	s.register()
	return s.Server.Serve(l)
}

// register 方法，用于注册服务到 etcd
func (s *Server) register() error {
	// 创建 etcd 客户端
	client, err := etcdv3.NewFromURL(s.EtcdAddr)
	if err != nil {
		return err
	}
	s.client = client
	// 创建 endpoints 管理器
	em, err := endpoints.NewManager(client, "service/"+s.Name)
	// 获取本机 IP 地址
	addr := netx.GetOutboundIP() + ":" + strconv.Itoa(s.Port)
	// 创建 etcd key
	key := "service/" + s.Name + "/" + addr

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 租期
	var ttl int64 = 5
	// 创建租约
	leaseResp, err := client.Grant(ctx, ttl)
	if err != nil {
		return err
	}

	// 添加 endpoint
	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		// 定位信息，客户端怎么连你
		Addr: addr,
	}, etcdv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}
	// 创建 keepalive 上下文
	kaCtx, kaCancel := context.WithCancel(context.Background())
	s.kaCancel = kaCancel
	// 创建 keepalive channel
	ch, err := client.KeepAlive(kaCtx, leaseResp.ID)
	// 启动 goroutine，处理 keepalive 响应
	go func() {
		//require.NoError(t, err1)
		for kaResp := range ch {
			// 记录日志
			s.L.Debug(kaResp.String())
		}
	}()
	return err
}

// Close 方法，用于关闭服务
func (s *Server) Close() error {
	// 取消 keepalive
	if s.kaCancel != nil {
		s.kaCancel()
	}
	// 关闭 etcd 客户端
	if s.client != nil {
		// 依赖注入，你就不要关
		return s.client.Close()
	}
	// 关闭 grpc 服务
	s.GracefulStop()
	return nil
}
