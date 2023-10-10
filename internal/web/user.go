package web

import (
	"fmt"
	"geekgo-webook/internal/domain"
	"geekgo-webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
)

type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	svc         *service.UserService
}

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

func NewUserHandler(svc *service.UserService) *UserHandler {
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)

	return &UserHandler{
		emailExp:    emailExp,
		passwordExp: passwordExp,
		svc:         svc,
	}
}

// users这个路由组下 注册了profile signup login edit等路由，每个HTTP请求和对应的处理方式HandlerFunc
func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users") // Group creates a new router group
	// type HandlerFunc func(*Context)
	ug.GET("/profile", u.Profile)
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit)
	ug.GET("/logout", u.Logout)
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	// profile 包括 Nickname ... 可以把数据放在结构体里 用ctx.JSON 来JSON化传给前端
	type Profile struct {
	}
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	// /signup请求 从客户端传了 邮箱密码 来发送注册请求 要用ctx接收前端发送的数据
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}

	// 从前端获取数据 存放在结构体里面之后 需要对邮箱格式和密码格式进行校验 regex
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字、特殊字符")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
	fmt.Printf("%v", req)
	// 注册成功 需要操作数据库 新增一条记录 可能用户已经注册过...
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	ctx.String(http.StatusOK, "注册成功")

	// todo 注册成功后 应该转去登录页面 或者 注册成功后直接设置成登录态

}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	// login 从前端接收 email password 和数据库中的进行比对
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	err := ctx.Bind(req)
	if err != nil {
		//Bind不成功会400
		return
	}
	// 接下来要和数据库里面的比对 调用repository层的方法 进行登录 登录成功之后 需要在服务器记录session
	// 调用Service的SignUp方法
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	fmt.Printf("%v", user)
	// todo 登录成功之后 创建jwt 使用jwt保持登录态 在middleware login_jwt.go中做jwt的登录态校验
	// jwt tokenstring 包含Header(加密算法) Payload(数据) Signature(签名)
	// 参考教程 https://pkg.go.dev/github.com/golang-jwt/jwt#example-New-Hmac
	token := jwt.New(jwt.SigningMethodHS512)
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
	fmt.Println(user)
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Login(ctx *gin.Context) {
	// login 从前端接收 email password 和数据库中的进行比对
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	err := ctx.Bind(req)
	if err != nil {
		//Bind不成功会400
		return
	}
	// 接下来要和数据库里面的比对 调用repository层的方法 进行登录 登录成功之后 需要在服务器记录session
	// 调用Service的SignUp方法
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	fmt.Printf("%v", user)
	// todo 登录成功之后 需要设置session 保持登录态
	sess := sessions.Default(ctx) // 返回Session interface Session具有 Get Set Delete Save等方法
	// 我可以随便设置值了
	// 你要放在 session 里面的值
	sess.Set("userId", user.Id)
	// 设置session参数  之后在登录校验中间件中添加session的刷新机制
	sess.Options(sessions.Options{
		Secure:   true, // 只能用Https协议
		HttpOnly: true, //
		MaxAge:   60,
	})
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")
	return

}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	// 我可以随便设置值了
	// 你要放在 session 里面的值
	sess.Options(sessions.Options{
		//Secure: true,
		//HttpOnly: true,
		MaxAge: -1,
	})
	sess.Save()
	ctx.String(http.StatusOK, "退出登录成功")
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	//
}
