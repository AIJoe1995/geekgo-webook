package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnknownForCode         = errors.New("我也不知发生什么了，反正是跟 code 有关")
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		client: client,
	}
}

// 编译器会在编译的时候，把 set_code 的代码放进来这个 luaSetCode 变量里
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

func (cache *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {

	// 主要逻辑 如果1min内发过验证码 需要返回发送验证码太频繁 主要逻辑是在lua脚本里实现的
	res, err := cache.client.Eval(ctx, luaSetCode, []string{cache.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		// 毫无问题
		return nil
	case -1:
		// 发送太频繁
		return ErrCodeSendTooMany
	//case -2:
	//	return
	default:
		// 系统错误
		return errors.New("系统错误")
	}
}
func (cache *RedisCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := cache.client.Eval(ctx, luaVerifyCode, []string{cache.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		// 正常来说，如果频繁出现这个错误，你就要告警，因为有人搞你
		return false, ErrCodeVerifyTooManyTimes
	case -2:
		return false, nil
		//default:
		//	return false, ErrUnknownForCode
	}
	return false, ErrUnknownForCode
}

func (cache *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
