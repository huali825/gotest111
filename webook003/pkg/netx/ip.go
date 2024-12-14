package netx

import "net"

// GetOutboundIP 方法获得对外发送消息的 IP 地址   别人连接你所用到的ip
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
