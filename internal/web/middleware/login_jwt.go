package middleware

import (
	"encoding/gob"
	ijwt "geekgo-webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// 登录校验

type LoginJWTMiddlewareBuilder struct {
	paths []string
	ijwt.Handler
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func NewLoginJWTMiddlewareBuilder(jwtHdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: jwtHdl,
	}
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用 Go 的方式编码解码
	gob.Register(time.Now())

	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if path == ctx.Request.URL.Path {
				return
			}
		}
		// 使用jwt-token进行登录校验 从ctx.Header中的Authorization里取出前端传来的JWTtoken
		//tokenHeader := ctx.GetHeader("Authorization")
		//if tokenHeader == "" {
		//	// 没登录
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		//segs := strings.Split(tokenHeader, " ")
		//if len(segs) != 2 {
		//	// 没登录，有人瞎搞
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		//tokenStr := segs[1]
		tokenStr := l.ExtractToken(ctx)
		claims := &ijwt.UserClaims{} // 要传入指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return ijwt.AtKey, nil
		})
		//token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		//	return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		//})
		if err != nil {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// err 为 nil，token 不为 nil
		if token == nil || !token.Valid || claims.Uid == 0 {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// jwt保存了user-agent 对比现在这个请求和登录时候的是不是同一个useragent发送的
		if claims.UserAgent != ctx.Request.UserAgent() {
			// 严重的安全问题
			// 你是要监控
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = l.CheckSession(ctx, claims.Ssid)
		if err != nil {
			// 要么 redis 有问题，要么已经退出登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 刷新jwt token // 设置长短token之后 不需要刷新token
		//now := time.Now()
		//if claims.ExpiresAt.Sub(now) < time.Second*50 {
		//	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
		//	tokenStr, err = token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
		//	if err != nil {
		//		// 记录日志
		//		log.Println("jwt 续约失败", err)
		//	}
		//	ctx.Header("x-jwt-token", tokenStr)
		//}
		ctx.Set("claims", claims) //??
		//ctx.Set("userId", claims.Uid)
	}
}
