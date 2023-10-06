package entity

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type SMSEntity interface {
	Store(ctx context.Context, tplId string, args []string, numbers []string) error
	FindRetryJob() ([]SMS, error)
	SetJobSuccess(ctx context.Context, id int64) error
}

type smsEntity struct {
	db *gorm.DB
}

// 存储一个待重试的任务
func (s *smsEntity) Store(ctx context.Context, tplId string, args []string, numbers []string) error {
	var sms SMS
	sms.Args = args
	sms.Numbers = numbers
	sms.TplId = tplId
	sms.Status = RETRY_NOACTIVE
	return s.db.WithContext(ctx).Create(&sms).Error
}

func (s *smsEntity) FindRetryJob() ([]SMS, error) {
	var sms []SMS
	err := s.db.Where("status = ?", RETRY_FAIL).Or("status = ?", RETRY_NOACTIVE).Find(&sms).Error
	return sms, err
}

func (s *smsEntity) SetJobSuccess(ctx context.Context, id int64) error {
	return s.db.Model(&SMS{}).WithContext(ctx).Where("id = ?", id).Update("status", RETRY_SUCCESS).Error
}

func NewSMSEntity(db *gorm.DB) SMSEntity {
	return &smsEntity{db: db}
}

const (
	RETRY_SUCCESS = iota
	RETRY_FAIL
	// 还未重试
	RETRY_NOACTIVE
)

// 再次尝试状态
type RetryStatus int

func (r *RetryStatus) Scan(value interface{}) error {
	val, ok := value.(int64)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal RetryStatus value:", value))
	}
	*r = RetryStatus(val)
	return nil
}

func (r RetryStatus) Value() (driver.Value, error) {
	return int64(r), nil
}

type SMS struct {
	Id        int64 `gorm:"primaryKey,autoIncrement"`
	TplId     string
	Args      []string `gorm:"-"`
	Numbers   []string `gorm:"-"`
	ArgsDB    string   `gorm:"column:args"`    // 实际存储在数据库的字段
	NumbersDB string   `gorm:"column:numbers"` // 实际存储在数据库的字段
	// 枚举类型
	Status RetryStatus `gorm:"type:int"`
}
