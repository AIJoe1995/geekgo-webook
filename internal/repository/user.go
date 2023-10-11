package repository

import (
	"context"
	"geekgo-webook/internal/domain"
	"geekgo-webook/internal/repository/cache"
	"geekgo-webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	// repository 先从缓存查找 缓存没有查找数据库并写回缓存
	u, err := repo.cache.Get(ctx, uid)
	if err == nil {
		return domain.User{}, err
	}
	ue, err := repo.dao.FindById(ctx, uid) // dao返回的是dao.User
	if err != nil {
		return domain.User{}, err
	}
	u = domain.User{
		Id:       ue.Id,
		Email:    ue.Email,
		Password: ue.Password,
	}

	go func() {
		err = repo.cache.Set(ctx, u)
		if err != nil {
			// 我这里怎么办？
			// 打日志，做监控
			//return domain.User{}, err
		}
	}()
	return u, err
}

// FindById  只有缓存中没找到数据的时候才去数据库查找 避免缓存崩溃 大量请求发到数据库
func (repo *UserRepository) FindByIdV1(ctx context.Context, uid int64) (domain.User, error) {
	u, err := repo.cache.Get(ctx, uid)
	switch err {
	case nil:
		return u, err
	case cache.ErrKeyNotExist:
		ue, err := repo.dao.FindById(ctx, uid)
		if err != nil {
			return domain.User{}, err
		}
		u = domain.User{
			Id:       ue.Id,
			Email:    ue.Email,
			Password: ue.Password,
		}

		go func() {
			err = repo.cache.Set(ctx, u)
			if err != nil {
				// 我这里怎么办？
				// 打日志，做监控
				//return domain.User{}, err
			}
		}()
		return u, err
	default:
		return domain.User{}, err
	}
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	// repo 交给dao层去寻找用户
	u, err := repo.dao.FindByEmail(ctx, email)
	// dao 层返回的是dao.User{} 应该装换成domain.User实体
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}
