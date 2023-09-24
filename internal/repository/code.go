package repository

import (
	"context"
	"webook/internal/repository/cache"
)

type CodeRepository interface {
	Send(ctx context.Context, biz string, phone string, inputCode string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

var (
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
	ErrUnknownForCode    = cache.ErrUnknownForCode
)

type cachedCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(c cache.CodeCache) CodeRepository {
	return &cachedCodeRepository{c}
}

func (r *cachedCodeRepository) Send(ctx context.Context, biz string, phone string, inputCode string) error {
	return r.cache.Set(ctx, biz, phone, inputCode)
}

func (r *cachedCodeRepository) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return r.cache.Verify(ctx, biz, phone, inputCode)
}
