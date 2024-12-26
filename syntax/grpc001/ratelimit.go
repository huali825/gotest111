/*
 * @Date: 2024年12月25日02:09:04
 * @LastEditors: TMH
 * @LastEditTime: 2024年12月25日02:09:08
 * @FilePath: /goworkwebook/syntax/grpc001/ratelimit.go
 * @Description: 限流算法之 计数器
 */

package grpc001

import (
	"context"
	"github.com/ecodeclub/ekit/queue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"sync/atomic"
	"time"
)

// CounterLimiter 计数器算法
type CounterLimiter struct {
	cnt       *atomic.Int32 // 当前请求数量
	threshold int32         // 系统的最大请求数量
}

func (c *CounterLimiter) NewServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {

		cnt := c.cnt.Add(1)
		defer func() {
			c.cnt.Add(-1)
		}()

		if cnt > c.threshold {
			return nil, status.Errorf(codes.ResourceExhausted, "限流")
		}
		return handler(ctx, req)
	}
}

//===========================//===========================//===========================//===========================//===========================//

// FixedWindowLimiter 固定窗口
type FixedWindowLimiter struct {
	// 固定窗口 确保每个窗口内的请求不超过一定的值
	window    time.Duration // 窗口大小 (时长)
	lastStart time.Time     // 上一个窗口的开始时间

	cnt       int // 当前窗口内的请求数量
	threshold int // 每个窗口内的最大请求数量

	lock sync.Mutex
}

func (l *FixedWindowLimiter) NewServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {

		l.lock.Lock()     // 获取锁，确保以下代码块是线程安全的
		now := time.Now() // 获取当前时间
		// 检查当前时间是否在最后一次开始时间加上窗口时间之前
		if now.Before(l.lastStart.Add(l.window)) {
			l.lastStart = now // 更新最后一次开始时间为当前时间
			l.cnt = 1         // 重置计数器为1
		}

		l.cnt++ // 计数器加1
		// 如果当前计数器小于等于阈值
		if l.cnt <= l.threshold {
			l.lock.Unlock() // 释放锁

			// 调用处理函数并返回结果
			res, err := handler(ctx, req)
			return res, err
		}
		return nil, status.Errorf(codes.ResourceExhausted, "限流")
	}

}

//===========================//===========================//===========================//===========================//===========================//

// 滑动窗口
type SlidingWindowLimiter struct {
	window    time.Duration                  // 窗口大小 (时长)
	queue     queue.PriorityQueue[time.Time] // 优先队列，用于存储每个请求的时间戳
	lock      sync.Mutex                     // 线程锁，用于保证线程安全
	threshold int                            // 阈值，表示在窗口时间内允许的最大请求数
}

// NewServerInterceptor 创建一个新的gRPC服务器拦截器，用于限流
func (l *SlidingWindowLimiter) NewServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {
		l.lock.Lock()

		// 获取当前时间
		now := time.Now()
		// 如果当前队列长度小于阈值，则允许请求通过
		if l.queue.Len() < l.threshold {
			// 将当前时间加入队列
			_ = l.queue.Enqueue(time.Now())
			// 解锁
			l.lock.Unlock()
			// 调用实际的请求处理函数
			return handler(ctx, req)
		}

		// 循环处理队列中的元素
		for {
			// 查看队列的第一个元素
			first, _ := l.queue.Peek()
			// 如果第一个元素的时间在窗口之外，则将其移出队列
			if first.Before(now.Add(-l.window)) {
				_, _ = l.queue.Dequeue()
			} else {
				// 否则跳出循环
				break
			}
		}

		// 再次检查队列长度是否小于阈值，如果是，则允许请求通过
		if l.queue.Len() < l.threshold {
			// 将当前时间加入队列
			_ = l.queue.Enqueue(time.Now())
			// 解锁
			l.lock.Unlock()
			// 调用实际的请求处理函数
			return handler(ctx, req)
		}

		// 解锁
		l.lock.Unlock()
		// 如果队列长度超过阈值，则返回资源耗尽错误
		return nil, status.Errorf(codes.ResourceExhausted, "限流")
	}
}

// ===========================//===========================//===========================//===========================//===========================//
// 令牌桶
type TokenBucketLimiter struct {
	//capacity int           // 令牌桶的容量

	interval  time.Duration // 令牌生成的时间间隔
	buckets   chan struct{} // 令牌桶，用于存储令牌
	closeCh   chan struct{} // 关闭通道，用于关闭令牌桶
	closeOnce sync.Once
}

func NewTokenbucketLimiter(interval time.Duration, capacity int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		interval: interval,
		buckets:  make(chan struct{}, capacity), // 令牌桶的容量为100
	}
}

func (l *TokenBucketLimiter) NewServerInterceptor() grpc.UnaryServerInterceptor {
	ticker := time.NewTicker(l.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				select {
				case l.buckets <- struct{}{}:
				default:
					// 如果令牌桶已满，则不生成新的令牌
				}
			case <-l.closeCh:
				return
			}

		}
	}()

	return func(ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {
		select {
		case <-l.buckets:
			return handler(ctx, req)
		default:
			return nil, status.Errorf(codes.ResourceExhausted, "限流")
		}

	}
}

func (l *TokenBucketLimiter) Close() error {
	l.closeOnce.Do(func() {
		close(l.closeCh)
	})
	return nil
}
