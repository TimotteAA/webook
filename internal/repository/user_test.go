package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	cachemocks "webook/internal/repository/cache/mock"
	"webook/internal/repository/entity"
	entitymocks "webook/internal/repository/entity/mock"

	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

func TestUserRepositoryFindById(t *testing.T) {
	now := time.Now().UnixMilli()
	testCases := []struct {
		name string
		// userhandler需要的实例
		mock func(ctrl *gomock.Controller) (entity.UserEntity, cache.UserCache)

		ctx    context.Context
		userId int64

		wantUser  domain.User
		wantError error
	}{
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (entity.UserEntity, cache.UserCache) {
				u := entitymocks.NewMockUserEntity(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				// 直接缓存命中
				c.EXPECT().Get(gomock.Any(), int64(22)).Return(domain.User{
					Id:        int64(22),
					Password:  "1234566",
					Phone:     "11111111111",
					Email:     "123@qq.com",
					CreatetAt: now,
				}, nil)
				return u, c

			},
			ctx:    context.Background(),
			userId: int64(22),
			// 预期输出
			wantUser: domain.User{
				Id:        22,
				Password:  "1234566",
				Phone:     "11111111111",
				Email:     "123@qq.com",
				CreatetAt: now,
			},
			wantError: nil,
		},
		{
			name: "缓存未命中",
			mock: func(ctrl *gomock.Controller) (entity.UserEntity, cache.UserCache) {
				u := entitymocks.NewMockUserEntity(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				// 先查缓存，未命中
				c.EXPECT().Get(gomock.Any(), int64(22)).Return(domain.User{}, cache.ErrKeyNotExist)
				//然后查entity
				u.EXPECT().FindById(gomock.Any(), int64(22)).Return(entity.User{
					Id: int64(22),
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Phone: sql.NullString{
						String: "11111111111",
						Valid:  true,
					},
					Password:   "1234566",
					CreateTime: now,
					UpdateTime: now,
				}, nil)
				// 写入缓存
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:        int64(22),
					Password:  "1234566",
					Phone:     "11111111111",
					Email:     "123@qq.com",
					CreatetAt: now,
				}).Return(nil)
				return u, c

			},
			ctx:    context.Background(),
			userId: int64(22),
			// 预期输出
			wantUser: domain.User{
				Id:        22,
				Password:  "1234566",
				Phone:     "11111111111",
				Email:     "123@qq.com",
				CreatetAt: now,
			},
			wantError: nil,
		},
		{
			name: "用户没找到",
			mock: func(ctrl *gomock.Controller) (entity.UserEntity, cache.UserCache) {
				u := entitymocks.NewMockUserEntity(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				// 直接没找到key
				c.EXPECT().Get(gomock.Any(), int64(22)).Return(domain.User{}, cache.ErrKeyNotExist)
				u.EXPECT().FindById(gomock.Any(), int64(22)).Return(entity.User{}, entity.ErrUserNotFound)
				return u, c

			},
			ctx:    context.Background(),
			userId: int64(22),
			// 预期输出
			wantUser:  domain.User{},
			wantError: entity.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		// 关闭
		defer ctrl.Finish()
		userEntity, userCache := tc.mock(ctrl)
		userService := NewUserRepository(userEntity, userCache)

		u, err := userService.FindById(tc.ctx, tc.userId)

		assert.Equal(t, err, tc.wantError)
		assert.Equal(t, u, tc.wantUser)
	}
}
