package repository

import (
	"context"
	"geekgo-webook/internal/repository/cache"
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func (repo *CodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	// 调用cache的Store方法
	return repo.cache.Set(ctx, biz, phone, code)

}

func (repo *CodeRepository) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, inputCode)
}
