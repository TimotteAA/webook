package repository

import (
	"context"
	"database/sql"
	"fmt"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/entity"
)

var ErrUserDuplicate = entity.ErrUserDuplciate
var ErrUserNotFound = entity.ErrUserNotFound

type UserRepository struct {
	entity *entity.UserEntity
	cache  *cache.UserCache
}

func NewUserRepository(entity *entity.UserEntity, cache *cache.UserCache) *UserRepository {
	return &UserRepository{entity: entity, cache: cache}
}

func (repo *UserRepository) Create(ctx context.Context, user domain.User) error {
	return repo.entity.Create(ctx, repo.domainToEntity(user))
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	ue, err := repo.entity.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(ue), nil
}

func (repo *UserRepository) FindById(ctx context.Context, userId int64) (domain.User, error) {
	// 先去缓存里面找
	u, err := repo.cache.Get(ctx, userId)
	fmt.Println("缓存中的结果 ", u)
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

	user := repo.entityToDomain(ue)

	//	写会缓存，忽视err
	_ = repo.cache.Set(ctx, user)
	return user, nil
}

func (repo *UserRepository) Update(ctx context.Context, userId int64, nickname string, description string, birthday int64) (domain.User, error) {
	ue, err := repo.entity.Update(ctx, userId, nickname, description, birthday)
	if err != nil {
		return domain.User{}, nil
	}
	user := repo.entityToDomain(ue)
	//	写入缓存

	_ = repo.cache.Set(ctx, user)
	return user, nil
}

func (repo *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	ue, err := repo.entity.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(ue), err
}

func (repo *UserRepository) entityToDomain(ue entity.User) domain.User {
	return domain.User{
		Id:          ue.Id,
		Email:       ue.Email.String,
		NickName:    ue.Nickname,
		Description: ue.Description,
		BirthDay:    ue.Birthday,
		Phone:       ue.Phone.String,
		Password:    ue.Password,
	}
}

func (repo *UserRepository) domainToEntity(ud domain.User) entity.User {
	return entity.User{
		Id:          ud.Id,
		Email:       sql.NullString{String: ud.Email, Valid: ud.Email != ""},
		Phone:       sql.NullString{String: ud.Phone, Valid: ud.Phone != ""},
		Nickname:    ud.NickName,
		Description: ud.Description,
		Birthday:    ud.BirthDay,
		Password:    ud.Password,
	}
}
