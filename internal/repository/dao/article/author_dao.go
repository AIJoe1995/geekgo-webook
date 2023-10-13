package article

import (
	"context"
	"geekgo-webook/internal/domain"
	"gorm.io/gorm"
)

type ArticleAuthorDAO interface {
	UpdateByID(ctx context.Context, art domain.Article) error
	Insert(ctx context.Context, artn Article) (int64, error)
}

func NewAuthorDAO(db *gorm.DB) ArticleAuthorDAO {
	panic("implement me")
}
