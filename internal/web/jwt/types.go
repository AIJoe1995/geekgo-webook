package jwt

//jwt目前需要提供的功能主要有
//1. 登录时需要根据uid设置jwt-token 长短token各一个，
//2. 在登录态校验的时候需要parsewithclaims jwttoken把token和claims拿出来， 进行逻辑校验
//3. 在前端访问refresh-token的路由时， 需要更新token
//4. 在调用退出登录接口的时候，需要把token claims里标识这次登录的uuid放在redis黑名单里， 登录态校验和访问refresh-token的时候都要检验一下携带的uuid在不在redis黑名单里 （有过期时间） 同时退出登录要把ctx里的长短token设置为空
//5. 检查token uuid是不是在redis黑名单里
//
//给jwt-token 定义接口 然后在需要用到jwt的地方 如UserHandler 和OAuth2WechatHandler的地方 组合jwt接口

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	SetLoginToken(ctx *gin.Context, uid int64) error
	SetJWTToken(ctx *gin.Context, uid int64, ssid string) error
	ClearToken(ctx *gin.Context) error
	CheckSession(ctx *gin.Context, ssid string) error
	ExtractToken(ctx *gin.Context) string
}

type UserClaims struct {
	Uid       int64
	Ssid      string
	UserAgent string
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	Uid  int64
	Ssid string
	jwt.RegisteredClaims
}
