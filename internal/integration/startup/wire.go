//go:build wireinject

package startup

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
		// 测试的db代替开发的db
		InitTestDB, InitRedis,
		ioc.InitSMSService,
		dao.NewUserDAO,
		cache.NewUserCache,
		cache.NewCodeCache,
		repository.NewUserRepository,
		repository.NewCodeRepository,
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,
		web.NewUserHandler,
		web.NewArticleHandler,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
		ijwt.NewRedisJWTHandler,
		//ioc.InitWechatService,
		InitPhantomWechatService,
		NewWechatHandlerConfig,
		web.NewOAuth2WechatHandler,
		InitLog,
	)
	return new(gin.Engine)
}

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)
var userSvcProvider = wire.NewSet(
	dao.NewUserDAO,
	cache.NewUserCache,
	repository.NewUserRepository,
	service.NewUserService)

// 提供InitArticleHandler 简单的依赖注入， 方便测试article
func InitArticleHandler() *web.ArticleHandler {
	wire.Build(thirdProvider,
		//dao.NewGORMArticleDAO,
		service.NewArticleService,
		web.NewArticleHandler,
		//repository.NewArticleRepository,
	)
	return &web.ArticleHandler{}
}

func InitUserSvc() service.UserService {
	wire.Build(thirdProvider, userSvcProvider)
	return service.NewUserService(nil)
}

func InitJwtHdl() ijwt.Handler {
	wire.Build(thirdProvider, ijwt.NewRedisJWTHandler)
	return ijwt.NewRedisJWTHandler(nil)
}
