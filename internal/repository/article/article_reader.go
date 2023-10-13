package article

import (
	"context"
	"geekgo-webook/internal/domain"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
}
