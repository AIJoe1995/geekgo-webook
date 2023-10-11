package repository

import (
	"context"
	"geekgo-webook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany        = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
	ErrUnknownForCode         = cache.ErrUnknownForCode
)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

// struct 里面cache 应该是接口类型 方便使用wire
type codeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(cache cache.CodeCache) CodeRepository {
	return &codeRepository{
		cache: cache,
	}
}

func (repo *codeRepository) Store(ctx context.Context, biz, phone, code string) error {
	// 调用cache的Store方法
	return repo.cache.Set(ctx, biz, phone, code)

}

func (repo *codeRepository) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, inputCode)
}
