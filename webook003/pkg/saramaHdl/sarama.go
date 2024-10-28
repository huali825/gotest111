package saramaHdl

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"goworkwebook/webook003/pkg/logger"
)

//	type ConsumerGroupHandler interface {
//		Setup(ConsumerGroupSession) error
//		Cleanup(ConsumerGroupSession) error
//		ConsumeClaim(ConsumerGroupSession, ConsumerGroupClaim) error
//	}
var _ sarama.ConsumerGroupHandler = &Handler[any]{}

// Handler 实现此接口
type Handler[T any] struct {
	l  logger.LoggerV1
	fn func(msg *sarama.ConsumerMessage, data T) error
}

// Setup 设置消费者组会话
func (h Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup 清理消费者组会话
func (h Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 方法用于处理从 Kafka 消费组中获取的消息
func (h Handler[T]) ConsumeClaim(
	// session 参数表示消费组会话
	session sarama.ConsumerGroupSession,
	// claim 参数表示消费组中的分区
	claim sarama.ConsumerGroupClaim) error {

	// 获取分区中的消息
	messages := claim.Messages()
	// 遍历消息
	for msg := range messages {
		// 定义一个变量 t，类型为 T
		var t T
		// 将消息的值解析为 t
		err := json.Unmarshal(msg.Value, &t)
		// 如果解析失败，记录错误并跳过该消息
		if err != nil {
			h.l.Error("json.Unmarshal failed",
				logger.String("topic", msg.Topic),
				logger.Int32("partition", msg.Partition),
				logger.Int64("offset", msg.Offset),
				logger.Error(err))
			continue
		}
		// 调用 Handler 的 fn 方法处理消息
		err = h.fn(msg, t)
		// 如果处理失败，记录错误
		if err != nil {
			h.l.Error("处理消息失败",
				logger.String("topic", msg.Topic),
				logger.Int32("partition", msg.Partition),
				logger.Int64("offset", msg.Offset),
				logger.Error(err))
		}
		// 标记消息已被处理
		session.MarkMessage(msg, "")
	}

	return nil
}
