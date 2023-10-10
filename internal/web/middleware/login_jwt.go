package middleware

import (
	"encoding/gob"
	"geekgo-webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

// 登录校验

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
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
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			// 没登录，有人瞎搞
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{} // 要传入指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
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

		// 刷新jwt token
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
			if err != nil {
				// 记录日志
				log.Println("jwt 续约失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}

	}
}
