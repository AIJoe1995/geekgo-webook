package sms

// 提供不同服务发送短信的接口

import "context"

type Service interface {
	Send(ctx context.Context, tpl string, args []string, numbers ...string) error
}
