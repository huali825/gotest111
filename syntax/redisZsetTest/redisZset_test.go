package redisZsetTest

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestRdsZSet(t *testing.T) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// 添加元素到zset
	err := rdb.ZAdd(ctx, "myzset", redis.Z{
		Score:  1.0,
		Member: "one",
	}).Err()
	if err != nil {
		panic(err)
	}

	err = rdb.ZAdd(ctx, "myzset", redis.Z{
		Score:  2.0,
		Member: "two",
	}).Err()
	if err != nil {
		panic(err)
	}

	// 获取zset中的元素
	vals, err := rdb.ZRangeWithScores(ctx, "myzset", 0, -1).Result()
	if err != nil {
		panic(err)
	}
	for _, val := range vals {
		fmt.Println(val.Member, val.Score)
	}

	// 删除zset中的元素
	err = rdb.ZRem(ctx, "myzset", "one").Err()
	if err != nil {
		panic(err)
	}

	// 获取zset中的元素数量
	count, err := rdb.ZCard(ctx, "myzset").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("zset length:", count)
}
