package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// 定义 Kafka 服务器地址
var addr = []string{"localhost:9094"}

// 测试同步生产者
func TestSyncProducer(t *testing.T) {
	// 创建 Kafka 配置
	cfg := sarama.NewConfig()
	// 设置生产者返回成功消息
	cfg.Producer.Return.Successes = true
	// 创建同步生产者
	producer, err := sarama.NewSyncProducer(addr, cfg)
	// 设置生产者分区器
	cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	//cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	//cfg.Producer.Partitioner = sarama.NewHashPartitioner
	//cfg.Producer.Partitioner = sarama.NewManualPartitioner
	//cfg.Producer.Partitioner = sarama.NewConsistentCRCHashPartitioner
	//cfg.Producer.Partitioner = sarama.NewCustomPartitioner()
	// 断言创建生产者时没有错误
	assert.NoError(t, err)
	// 发送 100 条消息
	for i := 0; i < 100; i++ {
		_, _, err = producer.SendMessage(&sarama.ProducerMessage{
			Topic: "test_topic",
			Value: sarama.StringEncoder("这是一条消息"),
			// 会在生产者和消费者之间传递的
			Headers: []sarama.RecordHeader{
				{
					Key:   []byte("key1"),
					Value: []byte("value1"),
				},
			},
			Metadata: "这是 metadata",
		})
	}
}

// 测试异步发送 生产者(用户)
func TestAsyncProducer(t *testing.T) {
	// 创建 Kafka 配置
	cfg := sarama.NewConfig()
	// 设置生产者返回成功消息
	cfg.Producer.Return.Successes = true
	// 设置生产者返回错误消息
	cfg.Producer.Return.Errors = true
	// 创建异步生产者
	producer, err := sarama.NewAsyncProducer(addr, cfg)
	// 断言创建生产者时没有错误
	assert.NoError(t, err)
	// 获取生产者输入通道
	msgs := producer.Input()
	// 发送一条消息
	msgs <- &sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("这是一条消息 tmh" + time.Now().String()),
		// 会在生产者和消费者之间传递的
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("key1"),
				Value: []byte("value1"),
			},
		},
		Metadata: "这是 metadata",
	}

	// 从生产者成功消息通道中读取消息
	select {
	case msg := <-producer.Successes():
		t.Log("发送成功", string(msg.Value.(sarama.StringEncoder)))
	// 从生产者错误消息通道中读取消息
	case err := <-producer.Errors():
		t.Log("发送失败", err.Err, err.Msg)
	}
}
