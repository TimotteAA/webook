package tencent

import (
	"context"
	"github.com/go-playground/assert/v2"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
	"testing"
)

func TestSender(t *testing.T) {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		t.Fatal()
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")

	c, err := sms.NewClient(common.NewCredential(secretId, secretKey),
		"ap-nanjing", profile.NewClientProfile())
	if err != nil {
		t.Fatal(err)
	}
	s := NewService(c, "1400792075", "慢慢学途公众号")

	testCases := []struct {
		name    string
		tplId   string
		params  []string
		numbers []string
		wantErr error
	}{
		{
			name:    "发送验证码",
			tplId:   "1675686",
			params:  []string{"123456", "5"},
			numbers: []string{"+8617301780942"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.Send(context.Background(), tc.tplId, tc.params, tc.numbers...)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
