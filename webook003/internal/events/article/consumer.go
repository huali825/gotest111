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
			saramaHdl.NewHandler[ReadEvent](i.l, i.Consume))
		if er != nil {
			i.l.Error("退出消费", logger.Error(er))
		}
	}()
	return err
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
