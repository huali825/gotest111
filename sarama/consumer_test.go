package sarama

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"log"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addr, "demo", cfg)
	assert.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()
	start := time.Now()
	err = consumer.Consume(ctx,
		[]string{"test_topic"}, ConsumerHandler{})
	assert.NoError(t, err)
	t.Log(time.Since(start).String())
}

type ConsumerHandler struct {
}

func (c ConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Println("这是 Setup")
	//partitions := session.Claims()["test_topic"]
	//for _, part := range partitions {
	//	session.ResetOffset("test_topic",
	//		part, sarama.OffsetOldest, "")
	//}
	return nil
}

func (c ConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("这是 Cleanup")
	return nil
}

func (c ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	// 获取消息
	msgs := claim.Messages()
	// 定义批次大小
	const batchSize = 10
	// 无限循环
	for t := 0; t < 5; t++ {
		// 打印一个批次开始
		log.Println("一个批次开始")
		// 创建一个批次
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)
		// 创建一个上下文，设置超时时间为5秒
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		// 定义一个done变量，用于控制循环
		var done = false
		// 创建一个errgroup，用于并发处理
		var eg errgroup.Group
		// 循环获取消息
		for i := 0; i < batchSize && !done; i++ {
			// 从消息通道中获取消息
			select {
			case <-ctx.Done():
				// 超时了
				done = true
			case msg, ok := <-msgs:
				if !ok {
					// 消息通道关闭了
					cancel()
					return nil
				}
				// 将消息添加到批次中
				batch = append(batch, msg)
				// 并发处理
				eg.Go(func() error {
					// 打印消息内容
					log.Println(string(msg.Value))
					return nil
				})
			}
		}
		// 取消上下文
		cancel()
		// 等待并发处理完成
		err := eg.Wait()
		if err != nil {
			// 打印错误
			log.Println(err)
			// 继续下一个批次
			continue
		}
		// 凑够了一批，然后你就处理
		// log.Println(batch)

		// 标记消息已处理
		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}
	return nil

}
func (c ConsumerHandler) ConsumeClaimV1(session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		log.Println(string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}
