package tencent

import (
	"context"
	"fmt"
	"webook/internal/repository"
)

type MemoryService struct {
	codeRepo *repository.CodeRepository
}

func NewMemoryService(repo *repository.CodeRepository) *MemoryService {
	return &MemoryService{
		codeRepo: repo,
	}
}

func (s *MemoryService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {

	fmt.Printf("验证码是 %s", args[0])
	return nil
}
