package ioc

import (
	"geekgo-webook/internal/service/sms"
	"geekgo-webook/internal/service/sms/memory"
	"geekgo-webook/internal/service/sms/ratelimit"
	"geekgo-webook/internal/service/sms/retryable"
	limiter "geekgo-webook/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitSMSService() sms.Service {
	return memory.NewService()
}

// // 初始化带有限流机制的短信服务
func InitSMSServiceLimit(cmd redis.Cmdable) sms.Service {
	// 换内存，还是换别的
	svc := ratelimit.NewRatelimitSMSService(memory.NewService(),
		limiter.NewRedisSlidingWindowLimiter(cmd, time.Second, 100))
	return retryable.NewService(svc, 3)
	return memory.NewService()
}
