package failover

import (
	"context"
	"errors"
	"time"
	"webook/internal/repository"
	"webook/internal/service/sms"
)

// 如果某次发送失败了，则将其记录在数据库中，异步的定时重做（可能短信服务不适用，但是逻辑应该是适用于大多数的第三方服务）
type failoverServiceV3 struct {
	svc     sms.Service
	smsRepo repository.SmsRepository
}

// 判断理由：如果一个第三方服务不可用，最直接会出现的情况是连续n次超时
func NewFailoverServiceV3(svc sms.Service, smsRepo repository.SmsRepository) sms.Service {
	// 在ioc创建实例后，就开始执行重试

	s := &failoverServiceV3{svc: svc, smsRepo: smsRepo}
	go s.retry()
	return s
}

func (f *failoverServiceV3) retry() {
	// 5分钟定时做一次？
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 执行重试逻辑
			jobs, err := f.smsRepo.FindRetryJob()
			if err != nil {
				// log here?
			}
			if len(jobs) > 0 {
				for _, job := range jobs {
					ctx := context.Background()
					job := job
					// todo1：限制go程数量
					// todo2：重试上限
					go func() {
						err := f.Send(ctx, job.TplId, job.Args, job.Numbers...)
						if err == nil {
							// 重置任务的状态为成功
							_ = f.smsRepo.SetJobSuccess(ctx, job.Id)
						}
					}()
				}
			}
			//case <-ctx.Done():
			// 当上下文被取消时，退出循环
			//return
		}
	}
}

func (f *failoverServiceV3) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 发送
	err := f.svc.Send(ctx, tplId, args, numbers...)
	// 仅针对限流后的服务、限流错误的进行保存
	if err == errors.New("短信发送接口被限流") || err == errors.New("短信发送服务接口限流异常") {
		_ = f.smsRepo.Store(ctx, tplId, args, numbers)
		return err
	}
	return err
}
