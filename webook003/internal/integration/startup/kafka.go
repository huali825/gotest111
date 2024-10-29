package startup

import (
	"github.com/IBM/sarama"
)

// InitSaramaClient 初始化Sarama客户端
func InitSaramaClient() sarama.Client {
	// 创建一个新的Sarama配置
	scfg := sarama.NewConfig()
	// 设置生产者返回成功标志
	scfg.Producer.Return.Successes = true
	// 使用配置创建一个新的Sarama客户端
	client, err := sarama.NewClient([]string{"localhost:9094"}, scfg)
	// 如果创建客户端失败，则抛出异常
	if err != nil {
		panic(err)
	}
	// 返回创建的客户端
	return client
}

// InitSyncProducer 使用Sarama客户端初始化同步生产者
func InitSyncProducer(c sarama.Client) sarama.SyncProducer {
	// 使用Sarama客户端创建一个新的同步生产者
	p, err := sarama.NewSyncProducerFromClient(c)
	// 如果创建生产者失败，则抛出异常
	if err != nil {
		panic(err)
	}
	// 返回创建的生产者
	return p
}
