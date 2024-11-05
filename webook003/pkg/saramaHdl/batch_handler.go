package saramaHdl

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"goworkwebook/webook003/pkg/logger"
	"log"
	"time"
)

var intnum01 int = 0

// BatchHandler 是一个用于处理批量消息的消费者组处理器
// 定义一个泛型结构体BatchHandler，用于处理批量的消息
type BatchHandler[T any] struct {
	// 定义一个函数，用于处理消息
	fn func(msgs []*sarama.ConsumerMessage, ts []T) error
	// 定义一个logger，用于记录日志
	l logger.LoggerV1
}

func NewBatchHandler[T any](l logger.LoggerV1, fn func(msgs []*sarama.ConsumerMessage, ts []T) error) *BatchHandler[T] {
	return &BatchHandler[T]{fn: fn, l: l}
}

func (b *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 获取消息
	msgs := claim.Messages()
	// 定义批次大小
	const batchSize = 10
	// 循环处理消息
	for {
		log.Println(intnum01, "一个批次开始")
		intnum01++
		// 定义批次
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)
		// 定义消息体
		ts := make([]T, 0, batchSize)
		// 设置超时时间
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		// 定义是否完成
		var done = false
		// 循环处理消息
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				// 超时了
				done = true
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return nil
				}
				// 将消息加入批次
				batch = append(batch, msg)
				// 反序列消息体
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					b.l.Error("反序列消息体失败",
						logger.String("topic", msg.Topic),
						logger.Int32("partition", msg.Partition),
						logger.Int64("offset", msg.Offset),
						logger.Error(err))
					continue
				}
				// 将消息体加入批次
				batch = append(batch, msg)
				ts = append(ts, t)
			}
		}
		cancel()
		// 凑够了一批，然后你就处理
		err := b.fn(batch, ts)
		if err != nil {
			b.l.Error("处理消息失败",
				// 把真个 msgs 都记录下来
				logger.Error(err))
		}
		// 标记消息已处理
		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}
}
