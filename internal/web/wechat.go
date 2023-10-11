package web

import (
	"errors"
	"fmt"
	"geekgo-webook/internal/service"
	"geekgo-webook/internal/service/oauth2/wechat"
	ijwt "geekgo-webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"time"
)

// 微信扫码登录功能
// 注册路由 微信请求扫码的路由 和扫完码跳转回来的路由 在第三方向微信发请求的时候带上redirecturl 跳转会回来redirecturl并带上code和state

type OAuth2WechatHandler struct {
	svc      wechat.Service // web调用service层
	cfg      WechatHandlerConfig
	stateKey []byte
	userSvc  service.UserService
	ijwt.Handler
}

type WechatHandlerConfig struct {
	Secure bool
	//StateKey
}

func NewOAuth2WechatHandler(svc wechat.Service,
	userSvc service.UserService,
	jwtHdl ijwt.Handler,
	cfg WechatHandlerConfig) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:      svc,
		userSvc:  userSvc,
		Handler:  jwtHdl,
		stateKey: []byte("95osj3fUD7foxmlYdDbncXz4VD2igvf1"),
		cfg:      cfg,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	// 调用wechat service 构造出要跳转的微信url
	state := uuid.New()
	url, err := h.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.String(http.StatusOK, "构造扫码登录URL失败")
		return
	}
	// 这里把state存放到cookie里 之后方便核对
	if err = h.setStateCookie(ctx, state); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})

}

func (h *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间，你预期中一个用户完成登录的时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	})
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		return err
	}
	ctx.SetCookie("jwt-state", tokenStr,
		600, "/oauth2/wechat/callback", // 这里的path含义？？？
		"", h.cfg.Secure, true)
	return nil
}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	// 校验一下我的 state
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		return fmt.Errorf("拿不到 state 的 cookie, %w", err)
	}

	var sc StateClaims
	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("token 已经过期了, %w", err)
	}

	if sc.State != state {
		return errors.New("state 不相等")
	}
	return nil
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	// callback这里需要取出code 然后构造请求url把code发给微信 从微信拿到access_token ....
	code := ctx.Query("code")
	// 在把code发给微信 构造携带code的请求之前 先验证state
	err := h.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "登录失败",
		})
		return
	}
	info, err := h.svc.VerifyCode(ctx, code) // 返回wechatinfo
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 这里怎么办？
	// 从 userService 里面拿 uid
	u, err := h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	// 把u.Id放到jwttoken里
	err = h.SetLoginToken(ctx, u.Id)

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

type StateClaims struct {
	State string
	jwt.RegisteredClaims
}
