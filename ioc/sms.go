package ioc

import (
	"geekgo-webook/internal/service/sms"
	"geekgo-webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	return memory.NewService()
}
