package service

import (
	"context"
	"fmt"
	"math/rand"
	"webook/internal/repository"
	"webook/internal/service/sms"
)

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context,
		biz string, phone string, code string) (bool, error)
}

var (
	ErrCodeSendTooMany   = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = repository.ErrCodeVerifyTooMany
	ErrUnknownForCode    = repository.ErrUnknownForCode
)

const loginCodeTplId = "1400792075"

type codeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(r repository.CodeRepository, sms sms.Service) CodeService {
	return &codeService{
		repo: r,
		sms:  sms,
	}
}

// 发送某个业务biz的验证码，限制发送频率
// redis存储某个业务发送短信的可以：phone_code:$biz:phone
// 1.如果key不存在，则发送
// 2.如果key存在，但无过期时间，则系统异常
// 3.key存在。过期时间是否还有9分组（一分钟都没过去），限制发送
// 4. 发送
// 由于在redis中存储code，故按照分层，需要cache和repo
func (c *codeService) Send(ctx context.Context, biz string, phone string) error {
	code := c.generateCode()
	err := c.repo.Send(ctx, biz, phone, code)
	if err != nil {
		// 存在问题
		return err
	}
	//	发送验证码
	err = c.sms.Send(ctx, loginCodeTplId, []string{code, "10分钟"}, phone)
	return err
}

func (c *codeService) Verify(ctx context.Context,
	biz string, phone string, code string) (bool, error) {
	ok, err := c.repo.Verify(ctx, biz, phone, code)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (c *codeService) generateCode() string {
	// 产生0-999999的随机数
	code := rand.Intn(1000000)
	// 不足6位，前面补0
	return fmt.Sprintf("%06d", code)
}
