package repository

import (
	"context"
	"webook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
	ErrUnknownForCode    = cache.ErrUnknownForCode
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{c}
}

func (r *CodeRepository) Send(ctx context.Context, biz string, phone string, inputCode string) error {
	return r.cache.Set(ctx, biz, phone, inputCode)
}

func (r *CodeRepository) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return r.cache.Verify(ctx, biz, phone, inputCode)
}
