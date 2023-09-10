package sms

import "context"

type Service interface {
	// ctx，模板id，发送参数，手机号
	Send(ctx context.Context, tplId string, args []string, numbers ...string) error
}
