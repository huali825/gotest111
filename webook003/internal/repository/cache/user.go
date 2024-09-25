package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"goworkwebook/webook003/internal/domain"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache struct {
	// 传单机 Redis 可以
	// 传 cluster 的 Redis 也可以
	client     redis.Cmdable
	expiration time.Duration
}

// NewUserCache
// A 用到了 B，B 一定是接口
// A 用到了 B，B 一定是 A 的字段
// A 用到了 B，A 绝对不初始化 B，而是外面注入
// expiration 1s, 1m
func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

// 如果没有数据，返回一个特定的 error
func (cache *UserCache) Get(ctx context.Context, id int64) (domain.DMUser, error) {
	key := cache.key(id)
	// 数据不存在，err = redis.Nil
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.DMUser{}, err
	}
	var u domain.DMUser
	err = json.Unmarshal(val, &u)
	//if err != nil {
	//	return domain.User{}, err
	//}
	//return u, nil
	return u, err
}

// 设置用户缓存
func (cache *UserCache) Set(ctx context.Context, u domain.DMUser) error {
	// 将用户信息序列化为JSON格式
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	// 生成缓存键
	key := cache.key(u.Id)
	// 将用户信息存入缓存
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}
func (cache *UserCache) key(id int64) string {
	// 生成缓存键 user:info:1 这种
	return fmt.Sprintf("user:info:%d", id)
}

// main 函数里面初始化好
//var RedisClient *redis.Client

//func GetUser(ctx context.Context, id int64) {
//	RedisClient.Get()
//}

//type UnifyCache interface {
//	Get(ctx context.Context, key string)
//	Set(ctx context.Context, key string, val any, expiration time.Duration)
//}
//
//
//type NewRedisCache() UnifyCache {
//
//}
