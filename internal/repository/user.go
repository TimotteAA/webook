package repository

import (
	"context"
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
