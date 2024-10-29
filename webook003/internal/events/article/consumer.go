package article

import (
	"context"
	"github.com/IBM/sarama"
	"goworkwebook/webook003/internal/repository"
	"goworkwebook/webook003/pkg/logger"
	"goworkwebook/webook003/pkg/saramaHdl"
	"time"
)

type InteractiveReadEventConsumer struct {
	repo   repository.InteractiveRepository
	client sarama.Client
	l      logger.LoggerV1
}

func NewInteractiveReadEventConsumer(repo repository.InteractiveRepository,
	client sarama.Client, l logger.LoggerV1) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{repo: repo, client: client, l: l}
}

func (i *InteractiveReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{TopicReadEvent},
			saramaHdl.NewBatchHandler[ReadEvent](i.l, i.BatchConsume)) //使用批量处理
		if er != nil {
			i.l.Error("退出消费", logger.Error(er))
		}
	}()
	return err
}

func (i *InteractiveReadEventConsumer) StartV1() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{TopicReadEvent},
			saramaHdl.NewHandler[ReadEvent](i.l, i.Consume)) //使用单条处理
		if er != nil {
			i.l.Error("退出消费", logger.Error(er))
		}
	}()
	return err
}

// BatchConsume 函数用于批量消费消息和事件
func (i *InteractiveReadEventConsumer) BatchConsume(
	// msgs 参数为 []*sarama.ConsumerMessage 类型的消息
	msgs []*sarama.ConsumerMessage,
	// events 参数为 []ReadEvent 类型的事件
	events []ReadEvent) error {

	// bizs 用于存储事件中的业务类型
	bizs := make([]string, 0, len(events))
	// bizIds 用于存储事件中的业务ID
	bizIds := make([]int64, 0, len(events))
	// 遍历事件
	for _, evt := range events {
		// 将事件中的业务类型添加到 bizs 中
		bizs = append(bizs, "article")
		// 将事件中的业务ID添加到 bizIds 中
		bizIds = append(bizIds, evt.Aid)
	}
	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 在函数结束时取消上下文
	defer cancel()
	// 调用 repo 的 BatchIncrReadCnt 方法，批量增加阅读次数
	return i.repo.BatchIncrReadCnt(ctx, bizs, bizIds)
}

// Consume 函数用于处理从Kafka中消费的消息
func (i *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage,
	event ReadEvent) error {
	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 在函数结束时取消上下文
	defer cancel()
	// 调用repo中的IncrReadCnt函数，增加文章的阅读次数
	return i.repo.IncrReadCnt(ctx, "article", event.Aid)
}