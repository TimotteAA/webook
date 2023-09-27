package service

import (
	"context"
	"testing"
	"webook/internal/domain"
	"webook/internal/repository"
	repomocks "webook/internal/repository/mock"

	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

func TestLogin(t *testing.T) {
		// 定义测试用例
		testCases := []struct{
			// 测试用例名称
			name string
	
			// userhandler需要的实例
			mock func(ctrl *gomock.Controller) (repository.UserRepository)
			
			// 模拟的方法入参数
			ctx context.Context
			user domain.User

			// 想要的结果和error
			wantUser domain.User
			wantError error
		}{
			{
				name: "用户登录",
				mock: func (ctrl *gomock.Controller) repository.UserRepository {
					userRepo := repomocks.NewMockUserRepository(ctrl);
					// Login里调用了repo的FindByEmail，返回一个user，和wantUser一样
					userRepo.EXPECT().FindByEmail(context.Background(), "123@qq.com").Return(domain.User{
						Id: 123,
						Email: "123@qq.com",
						Password: "$2a$10$s51GBcU20dkNUVTpUAQqpe6febjXkRYvhEwa5OkN5rU6rw2KTbNUi",
						CreatetAt: 1212512413213,
					}, nil)

					return userRepo;
				},
				ctx: context.Background(),
				// 邮箱密码登录
				user: domain.User{
					Email: "123@qq.com",
					Password: "hello#world123",
				},

				// 想要的结果和error
				wantUser: domain.User{
					Id: 123,
					Email: "123@qq.com",
					Password: "$2a$10$s51GBcU20dkNUVTpUAQqpe6febjXkRYvhEwa5OkN5rU6rw2KTbNUi",
					CreatetAt: 1212512413213,
				},
				wantError: nil,
			},
			{
				name: "邮箱不存在",
				mock: func (ctrl *gomock.Controller) repository.UserRepository {
					userRepo := repomocks.NewMockUserRepository(ctrl);
					// Login里调用了repo的FindByEmail，返回一个user，和wantUser一样
					userRepo.EXPECT().FindByEmail(context.Background(), "123@qq.com").Return(domain.User{
					}, repository.ErrUserNotFound)
					return userRepo;
				},
				ctx: context.Background(),
				// 邮箱密码登录
				user: domain.User{
					Email: "123@qq.com",
					Password: "hello#world123",
				},
				wantError: ErrEmailOrPassWrong,
			},
			{
				name: "密码错误",
				mock: func (ctrl *gomock.Controller) repository.UserRepository {
					userRepo := repomocks.NewMockUserRepository(ctrl);
					// Login里调用了repo的FindByEmail，返回一个user，和wantUser一样
					userRepo.EXPECT().FindByEmail(context.Background(), "123@qq.com").Return(domain.User{
						Email: "123@qq.com",
						// 错误的密码
						Password:	"$2a$10$s51GBcU20dkNUVTpUAQqpe6febjXkRYvhEwa5OkN5rU6rw2KTbNUi",
					}, nil)
					return userRepo;
				},
				ctx: context.Background(),
				// 输入错误的密码
				user: domain.User{
					Email: "123@qq.com",
					Password: "hello#",
				},

				// 想要的结果和error
				wantUser: domain.User{
					Id: 123,
					Email: "123@qq.com",
					Password: "$2a$10$s51GBcU20dkNUVTpUAQqpe6febjXkRYvhEwa5OkN5rU6rw2KTbNUi",
					CreatetAt: 1212512413213,
				},
				wantError: ErrEmailOrPassWrong,
			},
			{
				name: "用户未找到",
				mock: func (ctrl *gomock.Controller) repository.UserRepository {
					userRepo := repomocks.NewMockUserRepository(ctrl);
					// Login里调用了repo的FindByEmail，返回一个user，和wantUser一样
					userRepo.EXPECT().FindByEmail(context.Background(), "123@qq.com").Return(domain.User{
					}, repository.ErrUserNotFound)
					return userRepo;
				},
				ctx: context.Background(),
				// 输入错误的密码
				user: domain.User{
					Email: "123@qq.com",
					Password: "hello#world123",
				},
				wantError: ErrEmailOrPassWrong,
			},
		}
	
		// 运行测试用例
		for _, tc := range testCases {
			ctrl := gomock.NewController(t)
			// 关闭
			defer ctrl.Finish()
			userRepo := tc.mock(ctrl);
			userService := NewUserService(userRepo)
			
			_, err := userService.Login(tc.ctx, tc.user)

			assert.Equal(t, err, tc.wantError)
			// if (tc.wantUser != nil) {
			// 	assert.Equal(t, user, tc.wantUser) // 检查返回的用户是否与预期相符
			// }
		}
}