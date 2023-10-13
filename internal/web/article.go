package web

import (
	"geekgo-webook/internal/domain"
	"geekgo-webook/internal/service"
	ijwt "geekgo-webook/internal/web/jwt"
	"geekgo-webook/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 确保实现了handler接口
var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc    service.ArticleService
	logger logger.LoggerV1
}

func NewArticleHandler(svc service.ArticleService, logger logger.LoggerV1) *ArticleHandler {
	return &ArticleHandler{
		svc:    svc,
		logger: logger,
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	ag := server.Group("/articles")
	ag.POST("/edit", h.Edit)
}

// TDD 写了测试框架后 来写Edit的逻辑
func (h *ArticleHandler) Edit(ctx *gin.Context) {
	type ArticleReq struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var req ArticleReq
	if err := ctx.Bind(&req); err != nil { // 注意&req传指针 修改req
		return
	}

	c := ctx.MustGet("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		// 你可以考虑监控住这里
		//ctx.AbortWithStatus(http.StatusUnauthorized)
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.logger.Error("未发现用户的 session 信息")
		return
	}

	// 校验输入
	if req.Title == "" {
		ctx.String(http.StatusOK, "title不能为空")
	}

	// 调用svc 创建article
	//1. 定义articleservice接口 提供svc.Save的壳子 2. 向articlehandler 里注入articleservice
	// 3. 创建了注入了新的依赖 需要更改wire
	// 在service层面处理是新建还是修改
	id, err := h.svc.Save(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: claims.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 打日志？
		h.logger.Error("保存帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})

}
