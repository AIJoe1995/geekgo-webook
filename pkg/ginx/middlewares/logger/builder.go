package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"sync/atomic"
	"time"
)

// MiddlewareBuilder 注意点：
// 1. 小心日志内容过多。URL 可能很长，请求体，响应体都可能很大，你要考虑是不是完全输出到日志里面
// 2. 考虑 1 的问题，以及用户可能换用不同的日志框架，所以要有足够的灵活性
// 3. 考虑动态开关，结合监听配置文件，要小心并发安全

type MiddlewareBuilder struct {
	allowReqBody  *atomic.Bool // 同时读写的并发问题
	allowRespBody bool
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

func NewBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		loggerFunc: fn,
	}
}

// 链式调用 来修改allowReqBody allowRespBody参数
func (b *MiddlewareBuilder) AllowReqBody(ok bool) *MiddlewareBuilder {
	b.allowReqBody.Store(ok)
	return b
}

func (b *MiddlewareBuilder) AllowRespBody(ok bool) *MiddlewareBuilder {
	b.allowRespBody = ok
	return b
}

func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			Method: ctx.Request.Method,
			Url:    url,
		}
		if b.allowReqBody.Load() && ctx.Request.Body != nil {
			body, _ := ctx.GetRawData()                   // // GetRawData returns stream data.
			reader := io.NopCloser(bytes.NewReader(body)) // NopCloser returns a ReadCloser with a no-op Close method wrapping
			ctx.Request.Body = reader                     // Body  io.ReadCloser
			if len(body) > 1024 {
				body = body[:1024]
			}
			al.ReqBody = string(body)
		}
		// response 不能像request一样 用ctx.Request.Body
		// ctx 里面有Write属性 是ResponseWriter interface类型， 在ResponseWriter上定义了十几个方法
		// 可以自己实现Writer 装饰其中部分方法，在回写response前 先写到我们自己的数据结构里, 我们自己要操作ResponseWriter
		//在自己的结构体里面组合gin.ResponseWriter和用来接收数据的结构体*AccessLog
		if b.allowRespBody {

			ctx.Writer = responseWriter{
				al:             al,
				ResponseWriter: ctx.Writer,
			}
		}
		defer func() {
			al.Duration = time.Since(start).String()
			b.loggerFunc(ctx, al)
		}()
		// 执行到业务逻辑
		ctx.Next() // 为什么要调用ctx.Next()
		//// Next should be used only inside middleware.
		//// It executes the pending handlers in the chain inside the calling handler.
		//// See example in GitHub.
	}
}

type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

// 装饰器模式
func (w responseWriter) WriteHeader(statusCode int) {
	w.al.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w responseWriter) Write(data []byte) (int, error) {
	w.al.RespBody = string(data)
	return w.ResponseWriter.Write(data)
}

func (w responseWriter) WriteString(data string) (int, error) {
	w.al.RespBody = data
	return w.ResponseWriter.WriteString(data)
}

type AccessLog struct {
	// HTTP 请求的方法
	Method string
	// Url 整个请求 URL
	Url      string
	Duration string
	ReqBody  string
	RespBody string
	Status   int
}
