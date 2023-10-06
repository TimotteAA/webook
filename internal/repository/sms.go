package repository

import (
	"context"
	"webook/internal/repository/entity"
)

type SmsRepository interface {
	Store(ctx context.Context, tplId string, args []string, numbers []string) error
	FindRetryJob() ([]entity.SMS, error)
	SetJobSuccess(ctx context.Context, id int64) error
}

type smsRepository struct {
	e entity.SMSEntity
}

func NewSmsRepository(entity entity.SMSEntity) SmsRepository {
	return &smsRepository{e: entity}
}

func (s *smsRepository) Store(ctx context.Context, tplId string, args []string, numbers []string) error {
	return s.e.Store(ctx, tplId, args, numbers)
}

func (s *smsRepository) FindRetryJob() ([]entity.SMS, error) {
	return s.e.FindRetryJob()
}

func (s *smsRepository) SetJobSuccess(ctx context.Context, id int64) error {
	return s.e.SetJobSuccess(ctx, id)
}
