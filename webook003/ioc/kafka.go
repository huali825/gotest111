package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"goworkwebook/webook003/internal/events"
	"goworkwebook/webook003/internal/events/article"
)

// InitSaramaClient 初始化Sarama客户端
func InitSaramaClient() sarama.Client {
	// 定义配置结构体
	type Config struct {
		Addr []string `yaml:"addr"`
	}
	// 声明配置变量
	var cfg Config
	// 从配置文件中解析配置
	err := viper.UnmarshalKey("kafka", &cfg)
	// 如果解析失败，则抛出异常
	if err != nil {
		panic(err)
	}
	// 创建Sarama配置
	scfg := sarama.NewConfig()
	// 设置生产者返回成功标志
	scfg.Producer.Return.Successes = true
	// 创建Sarama客户端
	client, err := sarama.NewClient(cfg.Addr, scfg)
	// 如果创建失败，则抛出异常
	if err != nil {
		panic(err)
	}
	// 返回Sarama客户端
	return client
}

func InitSyncProducer(c sarama.Client) sarama.SyncProducer {
	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		panic(err)
	}
	return p
}

func InitConsumers(c1 *article.InteractiveReadEventConsumer) []events.Consumer {
	return []events.Consumer{c1}
}
