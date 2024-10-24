package startup

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// 初始化MongoDB数据库
func InitMongoDB() *mongo.Database {
	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// 在函数结束时取消上下文
	defer cancel()

	// 创建一个命令监视器
	monitor := &event.CommandMonitor{
		// 当命令开始时打印命令
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			fmt.Println(evt.Command)
		},
	}
	// 创建一个客户端选项，设置URI和监视器
	opts := options.Client().
		ApplyURI("mongodb://root:example@localhost:27017/").
		SetMonitor(monitor)
	// 连接MongoDB数据库
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		// 如果连接失败，则抛出异常
		panic(err)
	}
	// 返回数据库
	return client.Database("webook")
}
