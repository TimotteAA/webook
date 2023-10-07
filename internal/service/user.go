package service

import (
	"context"
	"errors"
	"webook/internal/domain"
	"webook/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SignUp(ctx context.Context, user domain.User) error
	Login(ctx context.Context, user domain.User) (domain.User, error)
	Edit(ctx context.Context, userId int64, nickname string, description string, birthday int64) (domain.User, error)
	FindOne(ctx context.Context, userId int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWeChat(ctx context.Context, userInfo domain.WeChatResult) (domain.User, error)
}

var ErrUserDuplicate = repository.ErrUserDuplicate
var ErrEmailOrPassWrong = errors.New("邮箱或密码错误")
var ErrUserNotFound = errors.New("用户不存在")

type userService struct {
	repo repository.UserRepository
}

// UserService工厂函数
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// handler的ctx先一路带下来
func (us *userService) SignUp(ctx context.Context, user domain.User) error {
	// 对密码加密，然后调用repo的insert方法

	// 加密后的密码
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}

	user.Password = string(hash)

	return us.repo.Create(ctx, user)
}

func (us *userService) Login(ctx context.Context, user domain.User) (domain.User, error) {
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

	//	比较密码，不太好测
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		return domain.User{}, ErrEmailOrPassWrong
	}
	// 返回结果
	return u, nil
}

// 编辑用户
func (uc *userService) Edit(ctx context.Context, userId int64, nickname string, description string, birthday int64) (domain.User, error) {
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

func (uc *userService) FindOne(ctx context.Context, userId int64) (domain.User, error) {
	user, err := uc.repo.FindById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (uc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := uc.repo.FindByPhone(ctx, phone)
	// 下面确保至少不是用户没找到的error，已经注册过了、或者别的原因
	if err != repository.ErrUserNotFound {
		return u, err
	}
	err = uc.repo.Create(ctx, domain.User{Phone: phone})
	// 有错误，但不是用户重复的错误，可能是系统错误
	if err != nil && err != repository.ErrUserDuplicate {
		return domain.User{}, err
	}
	// 上面两个error判断兜底了重复注册的err
	return uc.repo.FindByPhone(ctx, phone)
}

func (uc *userService) FindOrCreateByWeChat(ctx context.Context, userInfo domain.WeChatResult) (domain.User, error) {
	// 先查用户是否存在，也就是已经注册
	u, err := uc.repo.FindByWeChat(ctx, userInfo.OpenId)
	// openId唯一索引、可能已经注册过了、或者别的问题
	if err != repository.ErrUserDuplicate {
		return u, err
	}
	// 到这里没注册
	err = uc.repo.Create(ctx, domain.User{WeChatOpenId: userInfo.OpenId})
	if err != nil && err != repository.ErrUserDuplicate {
		return domain.User{}, err
	}
	return uc.repo.FindByWeChat(ctx, userInfo.OpenId)
}
