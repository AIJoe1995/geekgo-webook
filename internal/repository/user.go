package repository

import (
	"context"
	"database/sql"
	"geekgo-webook/internal/domain"
	"geekgo-webook/internal/repository/cache"
	"geekgo-webook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, uid int64) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByWechat(ctx context.Context, openID string) (domain.User, error)
}

////cannot use dao (variable of type dao.UserDAO) as *dao.GORMUserDAO value in argument to reposit
////ory.NewUserRepository: need type assertion
//type userRepository struct {
//	dao   *dao.GORMUserDAO
//	cache *cache.RedisUserCache
//}

type userRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &userRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *userRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.domainToEntity(u))
}

func (repo *userRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(u), nil
}

func (repo *userRepository) FindByWechat(ctx context.Context, openID string) (domain.User, error) {
	u, err := repo.dao.FindByWechat(ctx, openID)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(u), nil
}

func (repo *userRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	// repository 先从缓存查找 缓存没有查找数据库并写回缓存
	u, err := repo.cache.Get(ctx, uid)
	if err == nil {
		return domain.User{}, err
	}
	ue, err := repo.dao.FindById(ctx, uid) // dao返回的是dao.User
	if err != nil {
		return domain.User{}, err
	}
	u = repo.entityToDomain(ue)

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
func (repo *userRepository) FindByIdV1(ctx context.Context, uid int64) (domain.User, error) {
	u, err := repo.cache.Get(ctx, uid)
	switch err {
	case nil:
		return u, err
	case cache.ErrKeyNotExist:
		ue, err := repo.dao.FindById(ctx, uid)
		if err != nil {
			return domain.User{}, err
		}
		u = repo.entityToDomain(ue)

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

func (repo *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	// repo 交给dao层去寻找用户
	u, err := repo.dao.FindByEmail(ctx, email)
	// dao 层返回的是dao.User{} 应该装换成domain.User实体
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(u), nil
}

func (repo *userRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}

func (repo *userRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			// 我确实有手机号
			Valid: u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Ctime:    u.Ctime.UnixMilli(),
	}
}
