package limited

import (
	"context"
	"errors"
	"fmt"
	_sms "webook/internal/service/sms"
	"webook/pkg/limiter"
)

var errLimited = errors.New("短信发送接口被限流")

// 利用装饰器实现的limited service
type limitedSMSService struct {
	// 老的smsService
	sms _sms.Service
	// 限流器
	limiter limiter.Limiter
}

func NewLimitedSMSService(sms _sms.Service, limiter limiter.Limiter) _sms.Service {
	return &limitedSMSService{
		// 用于限流的腾讯云sms实现
		sms:     sms,
		limiter: limiter,
	}
}

func (l *limitedSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	key := l.key()
	limited, err := l.limiter.Limit(ctx, key)
	if err != nil {
		return fmt.Errorf("短信发送服务接口限流异常")
	}
	if limited {
		return errLimited
	}
	return l.sms.Send(ctx, tplId, args, numbers...)
}

func (l *limitedSMSService) key() string {
	return "phone_login_limit"
}
