package repository

import (
	"context"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/entity"
	"webook/internal/service/cache"
)

var ErrUserDuplciateEmail = entity.ErrUserDuplciateEmail
var ErrUserNotFound = entity.ErrUserNotFound

type UserRepository struct {
	entity *entity.UserEntity
	cache  *cache.UserCache
}

func NewUserRepository(entity *entity.UserEntity, cache *cache.UserCache) *UserRepository {
	return &UserRepository{entity: entity, cache: cache}
}

func (repo *UserRepository) Create(ctx context.Context, user domain.User) error {
	return repo.entity.Create(ctx, entity.User{
		Email:    user.Email,
		Password: user.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.entity.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (repo *UserRepository) FindById(ctx context.Context, userId int64) (domain.User, error) {
	// 先去缓存里面找
	u, err := repo.cache.Get(ctx, userId)
	// 用户存在
	if err == nil {
		return u, nil
	}
	// redis出错、或者没找到，查找数据库
	// 此处所有的压力都来到数据库，可能数据库会挂壁
	ue, err := repo.entity.FindById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	user := domain.User{
		Id:          ue.Id,
		Email:       ue.Email,
		NickName:    ue.Nickname,
		Description: ue.Description,
		BirthDay:    time.UnixMilli(ue.Birthday).Format("2006-01-02"),
	}

	//	写会缓存，忽视err
	_ = repo.cache.Set(ctx, user)
	return user, nil
}

func (repo *UserRepository) Update(ctx context.Context, userId int64, nickname string, description string, birthday int64) (entity.User, error) {
	return repo.entity.Update(ctx, userId, nickname, description, birthday)
}

func (repo *UserRepository) Detail(ctx context.Context, userId int64) (domain.User, error) {
	u, err := repo.entity.FindById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:          u.Id,
		Email:       u.Email,
		NickName:    u.Nickname,
		Description: u.Description,
		BirthDay:    time.UnixMilli(u.Birthday).Format("2006-01-02"),
	}, nil
}
