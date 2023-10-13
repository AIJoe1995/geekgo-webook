package web

import "github.com/gin-gonic/gin"

// 确保实现了handler接口
var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	ag := server.Group("/articles")
	ag.POST("/edit", h.Edit)
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {

}
