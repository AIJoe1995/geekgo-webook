package article

import (
	"context"
	"gorm.io/gorm"
)

type ArticleReaderDAO interface {
	Upsert(ctx context.Context, artn Article) error
}

// PublishArticle 这个代表的是线上表
type PublishArticle struct {
	Article
}

func NewReaderDAO(db *gorm.DB) ArticleReaderDAO {
	panic("implement me")
}
