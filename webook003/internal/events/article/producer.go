package article

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

const TopicReadEvent = "article_read"

// Producer 消息队列发出消息
type Producer interface {
	ProduceReadEvent(evt ReadEvent) error
}

type ReadEvent struct {
	Aid int64
	Uid int64
}

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

func NewSaramaSyncProducer(producer sarama.SyncProducer) Producer {
	return &SaramaSyncProducer{producer: producer}
}

// ProduceReadEvent 函数用于将ReadEvent事件发送到Kafka主题
func (s *SaramaSyncProducer) ProduceReadEvent(evt ReadEvent) error {
	// 将ReadEvent事件转换为JSON格式
	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	// 将JSON格式的ReadEvent事件发送到Kafka主题
	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicReadEvent,
		Value: sarama.StringEncoder(val),
	})
	return err
}
