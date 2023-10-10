package main

import (
	"geekgo-webook/internal/repository"
	"geekgo-webook/internal/repository/dao"
	"geekgo-webook/internal/service"
	"geekgo-webook/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
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

func main() {
	server := gin.Default()
	// 解决跨域问题 使用gin middleware server.Use(HandlerFunc)将HandlerFunc作用于全部路由
	server.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"*"}, // []string{"http://localhost:3030"}
		//AllowCrendentials 是否允许带上用户认证信息（比如 cookie）
		//AllowHeader：业务请求中可以带上的头。
		//AllowOriginFunc：哪些来源是允许的。
		//AllowMethods:     []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		//ExposeHeaders: []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 你的开发环境
				return true
			}
			return strings.Contains(origin, "yourcompany.com")

		},
		MaxAge: 12 * time.Hour,
	}))

	// 使用session middleware 可以提取session 使用session处理登录态的问题 在登陆成功之后把session保存起来 然后再设置登录校验的middleware来校验session
	// 代码示例 https://github.com/gin-contrib/sessions
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	server.Use(sessions.Sessions("mysession", store)) // middleware每次请求都会走这里，
	// sessions.Sessions返回的是HandlerFunc 会创建一个session结构体
	// s := &session{name, c.Request, store, nil, false, c.Writer}
	//		c.Set(DefaultKey, s) c是gin.Context
	// sessions的使用，

	// 校验登录 有些请求路径需要忽略 不经过校验 比如/users/login
	server.Use(func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/users/login" || ctx.Request.URL.Path == "/users/signup" {
			return
		}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	})

	db := initDB()
	dao := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(dao)
	svc := service.NewUserService(repo)

	u := web.NewUserHandler(svc)
	// server 处理GET请求 GET is a shortcut for router.Handle("GET", path, handlers)
	//server.GET(relativePath, HandlerFunc) 对于给定的请求路径，以HandlerFunc来处理请求返回响应
	// 示例 gin还有参数路由 通配符路由等
	// context核心职责 处理请求 返回响应 ctx.Bind() ctx.Query() ctx.String() ctx.JSON()
	//  r.GET("/ping", func(c *gin.Context) { // "users/:name" ctx.Param("name") "views/*.html" ctx.Param(".html")
	//    c.JSON(http.StatusOK, gin.H{
	//      "message": "pong",
	//    })
	//  })
	u.RegisterRoutes(server)
	server.Run(":8080") // 网络编程

}
