package tencent

import (
	"context"
	"fmt"
	"github.com/TimotteAA/gokit/utils"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms"
)

// A用到了B，由外部传入
type Service struct {
	client   *sms.Client
	appId    *string
	signName *string
}

func NewService(c *sms.Client, appId string, signName string) *Service {
	return &Service{
		client:   c,
		appId:    utils.ToPtr[string](appId),
		signName: utils.ToPtr(signName),
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 按照文档，实例化一个发送请求对象
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = s.appId
	request.SignName = s.signName

	// ctx往下传
	request.SetContext(ctx)
	request.TemplateId = utils.ToPtr(tplId)
	request.TemplateParamSet = common.StringPtrs(args)
	request.PhoneNumberSet = common.StringPtrs(numbers)

	// 发短信
	response, err := s.client.SendSms(request)
	if err != nil {
		return err
	}
	// 各个手机号的发送状态
	for _, status := range response.Response.SendStatusSet {
		// 指针判nil
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送失败，code：%v，原因：%v", status.Code, status.Message)
		}
	}
	return nil
}
