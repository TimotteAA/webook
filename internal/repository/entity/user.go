package entity

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

// 操作User表的entity
type UserEntity struct {
	db *gorm.DB
}

var (
	ErrUserDuplciateEmail = errors.New("邮箱冲突")
)

// UserEntity工厂函数
func NewUserEntity(db *gorm.DB) *UserEntity {
	return &UserEntity{db: db}
}

// 开始定义CRUD方法，不知道返回啥，先返回error
func (entity *UserEntity) Create(ctx context.Context, u User) error {
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
			return ErrUserDuplciateEmail
		}
	}
	return err
}

// user表结构
type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	// 为了便于处理时间，时间统一用UTC+0下的时间戳
	CreateTime int64
	UpdateTime int64
}
