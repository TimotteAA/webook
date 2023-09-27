package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"webook/internal/service"
	svcmocks "webook/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

func TestUserHandlerSignup(t *testing.T) {
	// 定义测试用例
	testCases := []struct{
		// 测试用例名称
		name string

		// userhandler需要的实例
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		// 构造请求
		reqBuilder func(t *testing.T) *http.Request


		// 预期响应码
		wantCode int
		// 预期body string
		wantResponse Result
	}{
		{
			name: "注册成功",
			mock: func (ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				// mock的userService
				userService := svcmocks.NewMockUserService(ctrl);
				// service方法预期的输入与输出，输入应该无所谓？主要是你想要啥输出
				userService.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)

				codeService := svcmocks.NewMockCodeService(ctrl);
				return userService, codeService
			},
			reqBuilder: func(t *testing.T) *http.Request {
				// 传给request的body
				body := bytes.NewBuffer([]byte(`{"email":"123@qq.com", "password": "123456aA!", "repassword": "123456aA!"}`))
				req, err := http.NewRequest(http.MethodPost, "/user/signup", body)
				// 设置header
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					log.Fatal(err)
				}
				return req;
			},
			wantCode: http.StatusOK,
			wantResponse: Result{Code: 0, Msg: "注册成功"},
		},
		{
			name: "body数据解析失败",
			mock: func (ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				// 不调用，返回nil即可
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				// 传给request的body
				body := bytes.NewBuffer([]byte(`{"email': '1111'}`))
				req, err := http.NewRequest(http.MethodPost, "/user/signup", body)
				// 设置header
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					log.Fatal(err)
				}
				return req;
			},
			wantCode: 400,
			wantResponse: Result{Code: 4, Msg: "Bind data error"},
		},
		{
			name: "邮箱格式错误",
			mock: func (ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				// 不调用，返回nil即可
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				// 传给request的body
				body := bytes.NewBuffer([]byte(`{"email":"123"}`))
				req, err := http.NewRequest(http.MethodPost, "/user/signup", body)
				// 设置header
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					log.Fatal(err)
				}
				return req;
			},
			wantCode: 200,
			wantResponse: Result{Code: 4, Msg: "邮箱格式不正确"},
		},
		{
			name: "密码格式不对",
			mock: func (ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				// 不调用，返回nil即可
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				// 传给request的body
				body := bytes.NewBuffer([]byte(`{"email":"123@qq.com", "password": "1234", "repassword": "123456aA!"}`))
				req, err := http.NewRequest(http.MethodPost, "/user/signup", body)
				// 设置header
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					log.Fatal(err)
				}
				return req;
			},
			wantCode: 200,
			wantResponse: Result{Code: 4, Msg: "密码必须大于8位，包含数字、特殊字符"},
		},
		{
			name: "两次密码不一致",
			mock: func (ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				// 不调用，返回nil即可
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				// 传给request的body
				body := bytes.NewBuffer([]byte(`{"email":"123@qq.com", "password": "123456aA!", "repassword": "123456aA"}`))
				req, err := http.NewRequest(http.MethodPost, "/user/signup", body)
				// 设置header
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					log.Fatal(err)
				}
				return req;
			},
			wantCode: 200,
			wantResponse: Result{Code: 4, Msg: "密码与确认密码不一致！"},
		},
		{
			name: "用户已注册",
			mock: func (ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				// mock的userService
				userService := svcmocks.NewMockUserService(ctrl);
				// 对于service的error，直接返回对应的error即可
				userService.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(service.ErrUserDuplicate)
				return userService, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				// 传给request的body
				body := bytes.NewBuffer([]byte(`{"email":"123@qq.com", "password": "123456aA!", "repassword": "123456aA!"}`))
				req, err := http.NewRequest(http.MethodPost, "/user/signup", body)
				// 设置header
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					log.Fatal(err)
				}
				return req;
			},
			wantCode: 200,
			wantResponse: Result{Code: 4, Msg: "邮箱已被注册"},
		},
		{
			name: "userService系统异常",
			mock: func (ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				// mock的userService
				userService := svcmocks.NewMockUserService(ctrl);
				// 对于service的error，直接返回对应的error即可
				userService.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("系统错误"))
				return userService, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				// 传给request的body
				body := bytes.NewBuffer([]byte(`{"email":"123@qq.com", "password": "123456aA!", "repassword": "123456aA!"}`))
				req, err := http.NewRequest(http.MethodPost, "/user/signup", body)
				// 设置header
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					log.Fatal(err)
				}
				return req;
			},
			wantCode: 200,
			wantResponse: Result{Code: 5, Msg: "系统错误"},
		},
	}

	// 运行测试用例
	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		// 关闭
		defer ctrl.Finish()
		// mock的服务
		userService, codeService := tc.mock(ctrl);
		// userhandler实例
		userHandler := NewUserHandler(userService, codeService);

		// 注册路由
		server := gin.Default();
		userHandler.RegisterRoutes(server)
		// req
		req := tc.reqBuilder(t)
		// response
		response := httptest.NewRecorder()
		// 运行
		server.ServeHTTP(response, req)
		// 断言结果
		assert.Equal(t, response.Code, tc.wantCode)
		// 解码response.Body到Result对象
		var respBody Result
		err := json.NewDecoder(response.Body).Decode(&respBody)
		if err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		// 断言响应体
		assert.Equal(t, respBody, tc.wantResponse)
	}
}