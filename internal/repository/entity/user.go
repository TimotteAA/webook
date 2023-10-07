package entity

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

type UserEntity interface {
	Create(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, userId int64) (User, error)
	Update(ctx context.Context, userId int64, nickname string, description string, birthday int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindByWeChat(ctx context.Context, openId string) (User, error)
}

var (
	ErrUserDuplciate = errors.New("用户已注册")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

// 操作User表的entity
type userEntity struct {
	db *gorm.DB
}

// UserEntity工厂函数
func NewUserEntity(db *gorm.DB) UserEntity {
	return &userEntity{db: db}
}

// 开始定义CRUD方法，不知道返回啥，先返回error
func (entity *userEntity) Create(ctx context.Context, u User) error {
	// 在此处处理时间，存毫秒
	now := time.Now().UnixMilli()
	u.CreateTime = now
	u.UpdateTime = now
	err := entity.db.WithContext(ctx).Create(&u).Error
	if mySqlErr, ok := err.(*mysql.MySQLError); ok {
		// 用的mysql数据库，断言成mysqlerror
		const uniqueConflictErrNum uint16 = 1062
		if mySqlErr.Number == uniqueConflictErrNum {
			// 唯一索引异常
			return ErrUserDuplciate
		}
	}
	return err
}

// 根据email查找用户
func (entity *userEntity) FindByEmail(ctx context.Context, email string) (User, error) {
	// 注意此处的User是表结构的User
	var u User
	result := entity.db.WithContext(ctx).Where("email = ?", email).First(&u)
	return u, result.Error
}

func (entity *userEntity) FindById(ctx context.Context, userId int64) (User, error) {
	var u User
	result := entity.db.WithContext(ctx).Where("id = ?", userId).First(&u)
	return u, result.Error
}

// 更新
func (entity *userEntity) Update(ctx context.Context, userId int64, nickname string, description string, birthday int64) (User, error) {
	var user User
	updateMap := make(map[string]interface{})
	if nickname != "" {
		updateMap["Nickname"] = nickname
	}
	if description != "" {
		updateMap["Description"] = description
	}

	updateMap["Birthday"] = birthday

	result := entity.db.WithContext(ctx).Model(&user).Where("id = ?", userId).Updates(updateMap)
	if result.Error != nil {
		return User{}, result.Error
	}
	// 更新完再查一下
	err := entity.db.WithContext(ctx).First(&user, userId).Error
	return user, err
}

func (entity *userEntity) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	result := entity.db.WithContext(ctx).Where("phone = ?", phone).First(&u)
	return u, result.Error
}

func (entity *userEntity) FindByWeChat(ctx context.Context, openId string) (User, error) {
	var u User
	result := entity.db.WithContext(ctx).Where("openId = ?", openId).First(&u)
	return u, result.Error
}

// user表结构
type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// NullString的scan方法从数据库中读取的值，转换成go中的值;
	Email    sql.NullString `gorm:"unique"`
	Password string

	Nickname    string
	Birthday    int64
	Description string `gorm:"size:350"`

	// 手机号
	Phone sql.NullString `gorm:"unique"`

	// 微信相关的字段：应用下的用户id
	WeChatOpenId  string `gorm:"type=varchar(1024),unique"`
	WeChatUnionId string `gorm:"type=varchar(1024)"`

	// 为了便于处理时间，时间统一用UTC+0下的时间戳
	CreateTime int64
	UpdateTime int64
}
