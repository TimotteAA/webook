package service

import "context"

type CodeService struct {
}

func NewCodeService() *CodeService {
	return &CodeService{}
}

// 发送某个业务biz的验证码，限制发送频率
// redis存储某个业务发送短信的可以：phone_code:$biz:phone
// 1.如果key不存在，则发送
// 2.如果key存在，但无过期时间，则系统异常
// 3.key存在。过期时间是否还有9分组（一分钟都没过去），限制发送
// 4. 发送
func (c *CodeService) Send(ctx context.Context, biz string, phone string) error {

}

func (C *CodeService) Verify(ctx context.Context,
	biz string, phone string, code string) (bool, error) {

}
