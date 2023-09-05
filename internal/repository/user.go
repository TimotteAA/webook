package repository

import (
	"context"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/entity"
)

var ErrUserDuplciateEmail = entity.ErrUserDuplciateEmail
var ErrUserNotFound = entity.ErrUserNotFound

type UserRepository struct {
	entity *entity.UserEntity
}

func NewUserRepository(entity *entity.UserEntity) *UserRepository {
	return &UserRepository{entity: entity}
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
	u, err := repo.entity.FindById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id: u.Id,
	}, nil
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
