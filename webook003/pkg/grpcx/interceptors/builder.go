package interceptors

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"net"
	"strings"
)

type Builder struct {
}

// PeerName 获取对端应用名称
func (b *Builder) PeerName(ctx context.Context) string {
	return b.grpcHeaderValue(ctx, "app")
}

// PeerIP 获取对端ip
func (b *Builder) PeerIP(ctx context.Context) string {
	// 如果在 ctx 里面传入。或者说客户端里面设置了，就直接用它设置的
	// 有些时候你经过网关之类的东西，就需要客户端主动设置，防止后面拿到网关的 IP
	clientIP := b.grpcHeaderValue(ctx, "client-ip")
	if clientIP != "" {
		return clientIP
	}

	// 从grpc里取对端ip
	pr, ok2 := peer.FromContext(ctx)
	if !ok2 {
		return ""
	}
	if pr.Addr == net.Addr(nil) {
		return ""
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	if len(addSlice) > 1 {
		return addSlice[0]
	}
	return ""
}

// grpcHeaderValue 从传入的上下文中获取gRPC请求头中的指定键的值
func (b *Builder) grpcHeaderValue(ctx context.Context, key string) string {
	// 检查传入的键是否为空，如果为空则直接返回空字符串
	if key == "" {
		return ""
	}

	// 从上下文中提取元数据（metadata），如果提取失败则返回空字符串
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	// 从元数据中获取指定键的值，并将多个值用分号（;）连接成一个字符串返回
	return strings.Join(md.Get(key), ";")
}
