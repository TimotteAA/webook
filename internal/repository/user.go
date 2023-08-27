package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/entity"
)

var ErrUserDuplciateEmail = entity.ErrUserDuplciateEmail

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
