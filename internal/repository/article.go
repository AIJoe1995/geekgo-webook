package repository

import (
	"context"
	"geekgo-webook/internal/domain"
	"geekgo-webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}

type CachedArticleRepository struct {
	dao dao.ArticleDAO
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}

func (c CachedArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateById(ctx, c.domainToEntity(art))
}

func (c CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	//调用dao.Create来创建文章
	id, err := c.dao.Insert(ctx, c.domainToEntity(art))
	return id, err
}

func (c CachedArticleRepository) domainToEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
