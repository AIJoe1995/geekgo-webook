package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

//
//jwt目前需要提供的功能主要有
//1. 登录时需要根据uid设置jwt-token 长短token各一个，
//2. 在登录态校验的时候需要parsewithclaims jwttoken把token和claims拿出来， 进行逻辑校验
//3. 在前端访问refresh-token的路由时， 需要更新token
//4. 在调用退出登录接口的时候，需要把token claims里标识这次登录的uuid放在redis黑名单里， 登录态校验和访问refresh-token的时候都要检验一下携带的uuid在不在redis黑名单里 （有过期时间） 同时退出登录要把ctx里的长短token设置为空
//5. 检查token uuid是不是在redis黑名单里
//
//给jwt-token 定义接口 然后在需要用到jwt的地方 如UserHandler 和OAuth2WechatHandler的地方 组合jwt接口

var (
	AtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
	RtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvfx")
)

type RedisJWTHandler struct {
	cmd redis.Cmdable
}

func (r RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := r.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = r.setRefreshToken(ctx, uid, ssid)
	return err

}

func (h *RedisJWTHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := RefreshClaims{
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid: uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(RtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

func (r RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       uid,
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (r RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")

	claims := ctx.MustGet("claims").(*UserClaims) // claims是在校验登录态的中间件里放进去的
	return r.cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", claims.Ssid),
		"", time.Hour*24*7).Err()
}

// 校验登录态的中间件 需要checksession
func (r RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	_, err := r.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	return err
}

func (r RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	// 我现在用 JWT 来校验
	tokenHeader := ctx.GetHeader("Authorization")
	//segs := strings.SplitN(tokenHeader, " ", 2)
	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd: cmd,
	}
}
