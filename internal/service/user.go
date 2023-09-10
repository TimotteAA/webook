package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"webook/internal/domain"
	"webook/internal/repository"
)

var ErrUserDuplciateEmail = repository.ErrUserDuplciateEmail
var ErrEmailOrPassWrong = errors.New("邮箱或密码错误")
var ErrUserNotFound = errors.New("用户不存在")

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

func (us *UserService) Login(ctx context.Context, user domain.User) (domain.User, error) {
	// 先根据Email查找用户
	u, err := us.repo.FindByEmail(ctx, user.Email)
	// 用户没找到
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrEmailOrPassWrong
	}
	// 其他错误
	if err != nil {
		return domain.User{}, err
	}

	//	比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		return domain.User{}, ErrEmailOrPassWrong
	}
	// 返回结果
	return u, nil
}

// 编辑用户
func (uc *UserService) Edit(ctx context.Context, userId int64, nickname string, description string, birthday int64) (domain.User, error) {
	//	先查找用户是否存在
	if _, err := uc.repo.FindById(ctx, userId); err != nil {
		return domain.User{}, err
	}

	//	更新用户数据
	user, err := uc.repo.Update(ctx, userId, nickname, description, birthday)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (uc *UserService) FindOne(ctx context.Context, userId int64) (domain.User, error) {
	user, err := uc.repo.FindById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
