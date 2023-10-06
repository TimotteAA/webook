package limited

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
	"testing"
	"webook/internal/service/sms"
	smsmocks "webook/internal/service/sms/mocks"
	"webook/pkg/limiter"
	limitermocks "webook/pkg/limiter/mocks"
)

func TestLimitedSMSService_Send(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter)

		// 输入，直接不写好了
		wantError error
	}{
		{
			name: "发送成功",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				smsService := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				smsService.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return smsService, l
			},
		},
		{
			name: "触发限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				smsService := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				return smsService, l
			},
			wantError: errors.New("短信发送接口被限流"),
		},
		{
			name: "限流器故障",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				smsService := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, fmt.Errorf("短信发送服务接口限流异常"))
				return smsService, l
			},
			wantError: errors.New("短信发送服务接口限流异常"),
		},
	}
	ctrl := gomock.NewController(t)
	// 关闭
	defer ctrl.Finish()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			smsService, l := tc.mock(ctrl)
			limitedSmsService := NewLimitedSMSService(smsService, l)
			err := limitedSmsService.Send(context.Background(), "1234", []string{""})
			assert.Equal(t, tc.wantError, err)
		})
	}
}
