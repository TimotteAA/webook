package failover

import (
	"context"
	"errors"
	"webook/internal/service/sms"
)

// 一个一个轮询
type failoverService struct {
	svcs []sms.Service
}

func NewFailoverService(svcs []sms.Service) sms.Service {
	return &failoverService{svcs: svcs}
}

func (f *failoverService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 一个一个轮询，直到成功
	for _, svc := range f.svcs {

		err := svc.Send(ctx, tplId, args, numbers...)
		// 某个服务成功了，直接返回
		if err != nil {
			return nil
		}
	}
	return errors.New("所有的短信服务都发送失败")
}
