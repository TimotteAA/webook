package limiter

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

type Limiter interface {
	Limit(ctx context.Context, key string) (bool, error)
}

//go:embed script/slide_window.lua
var luaScript string

type limiter struct {
	cmd redis.Cmdable
	// 限流窗口大小
	interval time.Duration
	// 窗口内请求大小
	rate int
}

func NewLimiter(c redis.Cmdable, interval time.Duration, rate int) Limiter {
	return &limiter{cmd: c, interval: interval, rate: rate}
}

func (l *limiter) Limit(ctx context.Context, key string) (bool, error) {
	return l.cmd.Eval(ctx, luaScript, []string{key}, l.interval.Milliseconds(), l.rate, time.Now().UnixMilli()).Bool()
}
