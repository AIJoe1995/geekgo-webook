package main

import (
	"geekgo-webook/internal/repository/dao"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:1234@tcp(localhost:3306)/webook"))
	if err != nil {
		// 我只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return redisClient
}

// 修改为面向接口编程， 从dao->repository->service->web 逐层修改， A调用B 那么B应该是接口，dao是最底层被调用，从dao开始修改
// 把各层的struct 先抽象出接口 然后NewStruct的时候 要返回接口，重新wire

func main() {
	server := InitWebServer()
	server.Run(":8080")
}

//
//func oldmain() {
//	server := gin.Default()
//	// 解决跨域问题 使用gin middleware server.Use(HandlerFunc)将HandlerFunc作用于全部路由
//	server.Use(cors.New(cors.Config{
//		//AllowOrigins:     []string{"*"}, // []string{"http://localhost:3030"}
//		//AllowCrendentials 是否允许带上用户认证信息（比如 cookie）
//		//AllowHeader：业务请求中可以带上的头。
//		//AllowOriginFunc：哪些来源是允许的。
//		//AllowMethods:     []string{"POST", "GET"},
//		AllowHeaders:  []string{"Content-Type", "Authorization"},
//		ExposeHeaders: []string{"x-jwt-token"}, // 带上这个前端才能拿到x-jwt-token 前端拿到token后一般放在Authorization Bearer里
//
//		AllowCredentials: true,
//		AllowOriginFunc: func(origin string) bool {
//			if strings.HasPrefix(origin, "http://localhost") {
//				// 你的开发环境
//				return true
//			}
//			return strings.Contains(origin, "yourcompany.com")
//
//		},
//		MaxAge: 12 * time.Hour,
//	}))
//
//	// 使用session middleware 可以提取session 使用session处理登录态的问题 在登陆成功之后把session保存起来 然后再设置登录校验的middleware来校验session
//	// 代码示例 https://github.com/gin-contrib/sessions
//	//store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
//	//server.Use(sessions.Sessions("mysession", store)) // middleware每次请求都会走这里，
//
//	//// sessions.Sessions返回的是HandlerFunc 会创建一个session结构体
//	// s := &session{name, c.Request, store, nil, false, c.Writer}
//	//		c.Set(DefaultKey, s) c是gin.Context
//	// sessions的使用，
//
//	// 校验登录 有些请求路径需要忽略 不经过校验 比如/users/login
//	//server.Use(middleware.NewLoginMiddlewareBuilder().
//	//	IgnorePaths("/users/login").
//	//	IgnorePaths("/users/signup").
//	//	Build())
//
//	// jwt登录校验
//	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
//		IgnorePaths("/users/login").
//		IgnorePaths("/users/signup").
//		Build())
//
//	db := initDB()
//	dao := dao.NewUserDAO(db)
//	client := initRedis()
//	usercache := cache.NewUserCache(client)
//	repo := repository.NewUserRepository(dao, usercache)
//	svc := service.NewUserService(repo)
//	smssvc := memory.NewService()
//	codecache := cache.NewCodeCache(client)
//	coderepo := repository.NewCodeRepository(codecache)
//	codesvc := service.NewCodeService(coderepo, smssvc)
//	u := web.NewUserHandler(svc, codesvc)
//	// server 处理GET请求 GET is a shortcut for router.Handle("GET", path, handlers)
//	//server.GET(relativePath, HandlerFunc) 对于给定的请求路径，以HandlerFunc来处理请求返回响应
//	// 示例 gin还有参数路由 通配符路由等
//	// context核心职责 处理请求 返回响应 ctx.Bind() ctx.Query() ctx.String() ctx.JSON()
//	//  r.GET("/ping", func(c *gin.Context) { // "users/:name" ctx.Param("name") "views/*.html" ctx.Param(".html")
//	//    c.JSON(http.StatusOK, gin.H{
//	//      "message": "pong",
//	//    })
//	//  })
//	u.RegisterRoutes(server)
//	server.Run(":8080") // 网络编程
//
//}
