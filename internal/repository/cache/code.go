package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany   = errors.New("发送登录验证码过于频繁")
	ErrCodeVerifyTooMany = errors.New("验证次数太多")
	ErrUnknownForCode    = errors.New("验证码系统错误")
)

// 编译时，把对应lua脚本的内容放到这个变量
//
//go:embed lua/set_code.lua
var setCodeScript string

//go:embed lua/verify_code.lua
var verifyCodeScript string

type CodeCache struct {
	//	存手机号的redis client
	client redis.Cmdable
}

func NewCodeCache(c redis.Cmdable) *CodeCache {
	return &CodeCache{client: c}
}

func (cache *CodeCache) Set(ctx context.Context,
	biz string, phone string, inputCode string) error {
	// 执行lua脚本
	result, err := cache.client.Eval(ctx, setCodeScript, []string{cache.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return err
	}
	switch result {
	case -2:
		return ErrUnknownForCode
	case -1:
		return ErrCodeSendTooMany
	case 0:
		return nil
	default:
		return ErrUnknownForCode
	}
}

func (cache *CodeCache) Verify(ctx context.Context,
	biz string, phone, inputCode string) (bool, error) {
	result, err := cache.client.Eval(ctx, verifyCodeScript, []string{cache.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch result {
	case 0:
		return true, nil
	case -1:
		return false, ErrCodeVerifyTooMany
	case -2:
		// 验证失败
		return false, nil
	}
	//	未知错误
	return false, ErrUnknownForCode
}

func (cache *CodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_login:%s:%s", biz, phone)
}
