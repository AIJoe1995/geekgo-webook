package service

import (
	"context"
	"geekgo-webook/internal/domain"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error) // 需要实现一个Article的domain
}

func NewArticleService() ArticleService {
	return &articeService{}
}

type articeService struct {
}

func (a articeService) Save(ctx context.Context, art domain.Article) (int64, error) {
	return 1, nil
}
