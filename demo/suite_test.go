package demo

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ArticleTestSuite struct {
	suite.Suite
	//server *gin.Engine
	//db     *gorm.DB
}

func (s *ArticleTestSuite) TestABC() {
	s.T().Log("hello，这是测试套件")
}

func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
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
