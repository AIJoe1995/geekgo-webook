//go:build wireinject

package main

import (
	"geekgo-webook/internal/repository"
	"geekgo-webook/internal/repository/cache"
	"geekgo-webook/internal/repository/dao"
	"geekgo-webook/internal/service"
	"geekgo-webook/internal/web"
	ijwt "geekgo-webook/internal/web/jwt"
	"geekgo-webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// wire依赖注入 步骤
// wire.Build() 注册各种初始化方法
// cmd运行wire命令
// 各种初始化方法 如 InitDB InitRedis等可以放在一个单独的ioc包里
// 需要调整代码结构 如将中间件抽出一个函数 初始化， 将server的初始化方法如注册路由等抽出一个函数

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,
		ioc.InitSMSService,
		dao.NewUserDAO,
		cache.NewUserCache,
		cache.NewCodeCache,
		repository.NewUserRepository,
		repository.NewCodeRepository,
		service.NewUserService,
		service.NewCodeService,
		web.NewUserHandler,
		ioc.InitWebServer,
		ioc.InitMiddlewares,
		ijwt.NewRedisJWTHandler,
		ioc.InitWechatService,
		ioc.NewWechatHandlerConfig,
		web.NewOAuth2WechatHandler,
	)
	return new(gin.Engine)
}
