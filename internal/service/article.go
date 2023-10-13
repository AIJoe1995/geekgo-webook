package service

import (
	"context"
	"geekgo-webook/internal/domain"
	"geekgo-webook/internal/repository"
	"github.com/gin-gonic/gin"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx *gin.Context, article domain.Article) (int64, error) // 需要实现一个Article的domain
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articeService{
		repo: repo,
	}
}

type articeService struct {
	repo repository.ArticleRepository
}

// 测试web的Publish方法时 会mock ArticleService 提供输入输出 所以暂时不需要implement
func (a articeService) Publish(ctx *gin.Context, article domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (a articeService) Save(ctx context.Context, art domain.Article) (int64, error) {
	// 区分新建和修改
	if art.Id > 0 {
		return art.Id, a.repo.Update(ctx, art)
	}
	return a.repo.Create(ctx, art)

}
