package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode string

	// 定义错误信息
	ErrCodeSendTooMany   = errors.New("发送太频繁")
	ErrCodeVerifyTooMany = errors.New("验证太频繁")
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

// CodeCache 结构体，用于缓存验证码
type RedisCodeCache struct {
	cmd redis.Cmdable
}

// NewCodeCache 创建一个新的 CodeCache 实例
func NewCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		cmd: cmd,
	}
}

// Set 设置验证码
func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	// 调用 redis 的 Eval 方法执行 lua 脚本
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		// 调用 redis 出了问题
		return err
	}
	switch res {
	case -2:
		// 验证码存在，但是没有过期时间
		return errors.New("验证码存在，但是没有过期时间")
	case -1:
		// 发送太频繁
		return ErrCodeSendTooMany
	default:
		// 设置成功
		return nil
	}
}

// Verify 验证验证码
func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	// 调用 redis 的 Eval 方法执行 lua 脚本
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		// 调用 redis 出了问题
		return false, err
	}
	switch res {
	case -2:
		// 验证码不存在
		return false, nil
	case -1:
		// 验证太频繁
		return false, ErrCodeVerifyTooMany
	default:
		// 验证成功
		return true, nil
	}
}

// key 生成验证码的 key
func (c *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
