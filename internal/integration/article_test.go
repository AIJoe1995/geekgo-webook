package integration

import (
	"bytes"
	"encoding/json"
	"geekgo-webook/internal/integration/startup"
	"geekgo-webook/internal/web"
	ijwt "geekgo-webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (s *ArticleTestSuite) SetupSuite() {
	s.server = gin.Default()
	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("claims", &ijwt.UserClaims{
			Uid: 123,
		})
	})
	s.db = startup.InitTestDB()
	artHdl := web.NewArticleHandler()
	artHdl.RegisterRoutes(s.server)

}

func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCases := []struct {
		name string
		// 集成测试准备数据
		before func(t *testing.T)
		// 集成测试验证数据
		after func(t *testing.T)

		// 预期中的输入
		art Article

		// HTTP 响应码
		wantCode int
		// 我希望 HTTP 响应，带上帖子的 ID
		wantRes Result[int64]
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//构造请求
			//执行请求
			// 验证结果
			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewBuffer(reqBody)) // 这里为什么用bytes.NewBuffer 不用bytes.NewReader
			require.NoError(t, err)
			// 数据是 JSON 格式
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			// 这就是 HTTP 请求进去 GIN 框架的入口。
			// 当你这样调用的时候，GIN 就会处理这个请求
			// 响应写回到 resp 里
			s.server.ServeHTTP(resp, req) // testsuite里面放了gin.Engine和gorm.DB
			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			var webRes Result[int64]
			// resp.Body是Body *bytes.Buffer类型
			//json.Unmarshal() //  Unmarshal(data []byte, v any) 这里能不能unmarshall
			err = json.NewDecoder(resp.Body).Decode(&webRes) // 反序列化把resp.Body的值存到结构体webRes中
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, webRes)

		})
	}
}

// 设计预期输入
type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

//// 预期输出 any在反序列化的时候会出现问题 修改为泛型实现
//type Result struct {
//	// 这个叫做业务错误码
//	Code int    `json:"code"`
//	Msg  string `json:"msg"`
//	Data any      `json:"data"`
//}

type Result[T any] struct {
	// 这个叫做业务错误码
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
