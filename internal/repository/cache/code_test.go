package cache

import (
	"context"
	"github.com/go-playground/assert/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/mock/gomock"
	"testing"
	"webook/internal/repository/cache/redismocks"
)

func TestCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) redis.Cmdable

		wantCode int64
		wantErr  error

		biz   string
		phone string
		code  string
	}{
		{
			name: "设置验证码成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				redisClient := redismocks.NewMockCmdable(ctrl)
				// 特殊的响应，返回的是int64
				redisResponse := redis.NewCmdResult(int64(0), nil)
				redisClient.EXPECT().Eval(gomock.Any(), setCodeScript, []string{"phone:login:11111111111"}, "123456").Return(redisResponse)
				return redisClient
			},
			wantCode: 0,
			wantErr:  nil,

			biz:   "login",
			phone: "11111111111",
			code:  "123456",
		},
		{
			name: "发送验证码次数过多",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				redisClient := redismocks.NewMockCmdable(ctrl)
				// 特殊的响应，返回的是int64
				redisResponse := redis.NewCmdResult(int64(-1), nil)
				redisClient.EXPECT().Eval(gomock.Any(), setCodeScript, []string{"phone:login:11111111111"}, "123456").Return(redisResponse)
				return redisClient
			},
			wantCode: -1,
			wantErr:  ErrCodeSendTooMany,

			biz:   "login",
			phone: "11111111111",
			code:  "123456",
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				redisClient := redismocks.NewMockCmdable(ctrl)
				// 特殊的响应，返回的是int64
				redisResponse := redis.NewCmdResult(int64(-2), nil)
				redisClient.EXPECT().Eval(gomock.Any(), setCodeScript, []string{"phone:login:11111111111"}, "123456").Return(redisResponse)
				return redisClient
			},
			wantCode: -2,
			wantErr:  ErrUnknownForCode,

			biz:   "login",
			phone: "11111111111",
			code:  "123456",
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		codeCache := NewCodeCache(tc.mock(ctrl))
		ctx := context.Background()
		err := codeCache.Set(ctx, tc.biz, tc.phone, tc.code)
		assert.Equal(t, tc.wantErr, err)
	}
}
