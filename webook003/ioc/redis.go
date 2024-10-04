package ioc

import (
	"github.com/redis/go-redis/v9"
	"goworkwebook/webook003/config"
)

func InitRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
}
