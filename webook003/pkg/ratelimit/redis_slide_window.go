package ratelimit

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var luaSlideWindow string

// RedisSlidingWindowLimiter Redis 上的滑动窗口算法限流器实现
type RedisSlidingWindowLimiter struct {
	cmd redis.Cmdable

	// 窗口大小
	interval time.Duration

	// 阈值
	rate int
	// interval 内允许 rate 个请求

	// 1s 内允许 3000 个请求
}

func NewRedisSlidingWindowLimiter(cmd redis.Cmdable,
	interval time.Duration, rate int) Limiter {
	return &RedisSlidingWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

// Limit 是 RedisSlidingWindowLimiter 结构体的一个方法，用于执行滑动窗口限流逻辑
func (r *RedisSlidingWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	// 调用 r.cmd.Eval 方法执行 Lua 脚本 luaSlideWindow
	// ctx 是上下文对象，用于控制请求的取消和超时
	// luaSlideWindow 是 Lua 脚本，用于在 Redis 中实现滑动窗口限流算法
	// []string{key} 是 Lua 脚本中使用的键，这里传入了一个键，即限流的标识
	// r.interval.Milliseconds() 将限流的时间间隔转换为毫秒，作为 Lua 脚本的第一个参数
	// r.rate 是限流的速率，即时间窗口内允许的最大请求数，作为 Lua 脚本的第二个参数
	// time.Now().UnixMilli() 获取当前时间的毫秒值，作为 Lua 脚本的第三个参数
	// Eval 方法执行 Lua 脚本并返回结果，这里调用 .Bool() 方法将结果转换为布尔值
	// 返回值 (bool, error) 表示是否超过限流以及可能的错误信息
	return r.cmd.Eval(ctx, luaSlideWindow, []string{key},
		r.interval.Milliseconds(), r.rate, time.Now().UnixMilli()).Bool()
}
