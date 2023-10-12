package logger

import (
	"context"
	"geekgo-webook/internal/service/sms"
	"go.uber.org/zap"
)

type Service struct {
	svc sms.Service
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	zap.L().Debug("发送短信", zap.String("tpl", tpl),
		zap.Any("args", args))
	err := s.svc.Send(ctx, tpl, args, numbers...)
	if err != nil {
		zap.L().Debug("发送短信出现异常", zap.Error(err))
	}
	return err
}
