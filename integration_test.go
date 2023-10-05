package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"webook/internal/web"
	"webook/ioc"
)

// 对于完整的SERVER进行测试
func TestUserSendCode(t *testing.T) {
	sendSignUpCodeUrl := "/user/signup/code/send"
	// 利用wire组装的server
	server := InitWebServer()
	// redis
	rdb := ioc.InitRedis()

	testCases := []struct {
		name string
		// 准备数据的函数
		before func(t *testing.T)
		// 测试完成清理数据
		after func(t *testing.T)

		// 发送的手机号
		phone string

		// 预期response
		wantCode   int
		wantResult web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
				// 发送成功啥都不做
			},
			after: func(t *testing.T) {
				ctx := context.Background()
				key := "phone_code:login:22222"

				ttl, err := rdb.TTL(ctx, key).Result()
				require.NoError(t, err)
				println("ttl length ", ttl)
				//require.True(t, ttl > 9*time.Minute)

				val, err := rdb.GetDel(ctx, key).Result()
				println("val ", val)
				require.NoError(t, err)
				require.True(t, len(val) == 6)
			},
			phone:    "22222",
			wantCode: 200,
			wantResult: web.Result{
				Code: 0,
				Msg:  "发送成功",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 准备数据
			tc.before(t)

			// 请求body
			reqBody := fmt.Sprintf(`{"phone": "%s"}`, tc.phone)
			// 制造请求
			req, err := http.NewRequest(http.MethodPost, sendSignUpCodeUrl, bytes.NewBuffer([]byte(reqBody)))
			require.NoError(t, err)
			// response
			response := httptest.NewRecorder()
			// 运行
			server.ServeHTTP(response, req)

			code := response.Code
			res, _ := rdb.Get(context.Background(), "phone_code:login:921").Result()
			println("result ", res)
			println(code)
			// 反序列化结果
			var result web.Result
			err = json.Unmarshal(response.Body.Bytes(), &result)
			//println(result)
			require.NoError(t, err)
			assert.Equal(t, code, tc.wantCode)
			assert.Equal(t, result, tc.wantResult)
			// 清理数据
			tc.after(t)
		})
	}
}
