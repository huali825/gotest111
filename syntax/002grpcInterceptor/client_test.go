package grpcInterceptor

import (
	"context"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
