package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProducerMain(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(addr, cfg)
	cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	//cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	//cfg.Producer.Partitioner = sarama.NewHashPartitioner
	//cfg.Producer.Partitioner = sarama.NewManualPartitioner
	//cfg.Producer.Partitioner = sarama.NewConsistentCRCHashPartitioner
	//cfg.Producer.Partitioner = sarama.NewCustomPartitioner()
	assert.NoError(t, err)
	for i := 0; i < 10000; i++ {
		_, _, err = producer.SendMessage(&sarama.ProducerMessage{
			Topic: "article_read",
			Value: sarama.StringEncoder(`{"aid": 1, "uid": 123}`),
			// 会在生产者和消费者之间传递的
			//Headers: []sarama.RecordHeader{
			//	{
			//		Key:   []byte("key1"),
			//		Value: []byte("value1"),
			//	},
			//},
			//Metadata: "这是 metadata",
		})
	}
}
