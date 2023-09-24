package local

import (
	"context"
	"fmt"
)

type MemoryService struct {
}

func NewMemoryService() *MemoryService {
	return &MemoryService{}
}

func (s *MemoryService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {

	fmt.Printf("验证码是 %s\n", args[0])
	return nil
}
