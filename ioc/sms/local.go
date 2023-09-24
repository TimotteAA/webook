package sms

import (
	"webook/internal/service/sms"
	"webook/internal/service/sms/local"
)

func InitSmsService() sms.Service {
	return local.NewMemoryService()
}
