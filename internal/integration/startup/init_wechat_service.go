package startup

import (
	"geekgo-webook/internal/service/oauth2/wechat"
)

// InitPhantomWechatService 没啥用的虚拟的 wechatService

func InitPhantomWechatService() wechat.Service {
	return wechat.NewService("", "")
}
