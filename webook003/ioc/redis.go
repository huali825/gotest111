package ioc

import (
	"context"
	"github.com/redis/go-redis/v9"
	"goworkwebook/webook003/config"
)

//	func InitRedis() redis.Cmdable {
//		return redis.NewClient(&redis.Options{
//			Addr: config.Config.Redis.Addr,
//		})
//	}

func InitRedis() redis.Cmdable {
	client := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic("Redis connection failed: " + err.Error())
	}

	return client
}
