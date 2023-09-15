package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// 编译时，把对应lua脚本的内容放到这个变量
//
//go:embed lua/set_code.lua
var setCodeScript string

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
		return errors.New("系统错误")
	case -1:
		return errors.New("发送登录验证码过于频繁")
	case 0:
		return nil
	default:
		return errors.New("系统错误")
	}
}

func (cache *CodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_login:%s:%s", biz, phone)
}
