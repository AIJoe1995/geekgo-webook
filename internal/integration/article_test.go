package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"geekgo-webook/internal/integration/startup"
	"geekgo-webook/internal/repository/dao"
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
	artHdl := startup.InitArticleHandler()
	artHdl.RegisterRoutes(s.server)

}

// 每一个测试运行完之后都会执行
func (s *ArticleTestSuite) TearDownTest() {
	s.db.Exec("truncate Table articles")
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
		{
			name: "新建帖子-保存成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				// 验证数据库

				// 从数据库中拿到dao.Article 来对比数据
				var art dao.Article
				err := s.db.Where("id = ?", 1).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       1,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
				}, art)

			},
			art: Article{
				Title:   "我的标题",
				Content: "我的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 1,
				Msg:  "OK",
			},
		},
		{
			name: "修改帖子-保存成功",
			before: func(t *testing.T) {
				// 准备库里的帖子 修改这个帖子
				err := s.db.Create(dao.Article{
					Id:       2,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
					// 跟时间有关的测试，不是逼不得已，不要用 time.Now()
					// 因为 time.Now() 每次运行都不同，你很难断言
					Ctime: 123,
					Utime: 234,
				}).Error
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				// 验证数据库

				// 从数据库中拿到dao.Article 来对比数据
				var art dao.Article
				err := s.db.Where("id = ?", 2).First(&art).Error
				assert.NoError(t, err)

				assert.True(t, art.Utime > 234)

				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       2,
					Title:    "新的标题",
					Content:  "新的内容",
					AuthorId: 123,
					Ctime:    123,
				}, art)

			},
			art: Article{
				Id:      2,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 2,
				Msg:  "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//构造请求
			//执行请求
			// 验证结果
			fmt.Println("执行tc.before开始")
			tc.before(t)
			fmt.Println("执行tc.before结束")
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
			fmt.Println("执行tc.after开始")
			tc.after(t)
			fmt.Println("执行tc.after结束")
		})
	}
}

func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

// 设计预期输入
type Article struct {
	Id      int64
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
