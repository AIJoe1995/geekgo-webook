package ioc

import (
	"geekgo-webook/internal/service/oauth2/wechat"
	"geekgo-webook/internal/web"
)

func InitWechatService() wechat.Service {
	//appId, ok := os.LookupEnv("WECHAT_APP_ID")
	//if !ok {
	//	panic("没有找到环境变量 WECHAT_APP_ID ")
	//}
	//appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	//if !ok {
	//	panic("没有找到环境变量 WECHAT_APP_SECRET")
	//}
	//
	appId, appKey := "1", "692jdHsogrsYqxaUK9fgxw"
	// 692jdHsogrsYqxaUK9fgxw
	return wechat.NewService(appId, appKey)
}

func NewWechatHandlerConfig() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secure: false,
	}
}
