package service

import (
	"context"
	"geekgo-webook/internal/domain"
	"geekgo-webook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error) // 需要实现一个Article的domain
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articeService{
		repo: repo,
	}
}

type articeService struct {
	repo repository.ArticleRepository
}

func (a articeService) Save(ctx context.Context, art domain.Article) (int64, error) {
	id, err := a.repo.Create(ctx, art)
	return id, err

}
