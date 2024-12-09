package etcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"testing"
	"time"
)

const (
	serviceName = "my-service"
	serviceAddr = "127.0.0.1:8080"
	etcdAddr    = "127.0.0.1:2379"
)

func TestSunafa11111(t *testing.T) {
	// 创建etcd客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(cli)

	// 创建租约
	resp, err := cli.Grant(context.Background(), 10)
	if err != nil {
		log.Fatal(err)
	}

	// 注册服务
	key := serviceName + "/" + serviceAddr
	_, err = cli.Put(context.Background(), key, serviceAddr, clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}

	// 续约
	_, err = cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		log.Fatal(err)
	}

	// 服务发现
	getResp, err := cli.Get(context.Background(), serviceName+"/", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	for _, kv := range getResp.Kvs {
		log.Printf("Found service: %s", string(kv.Value))
	}

	// 退出时注销服务
	defer func() {
		_, err = cli.Revoke(context.Background(), resp.ID)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// 保持程序运行
	select {}
}
