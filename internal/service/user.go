package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"webook/internal/domain"
	"webook/internal/repository"
)

var ErrUserDuplciateEmail = repository.ErrUserDuplciateEmail

type UserService struct {
	repo *repository.UserRepository
}

// UserService工厂函数
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// handler的ctx先一路带下来
func (us *UserService) SignUp(ctx context.Context, user domain.User) error {
	// 对密码加密，然后调用repo的insert方法

	// 加密后的密码
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}

	user.Password = string(hash)

	return us.repo.Create(ctx, user)
}
