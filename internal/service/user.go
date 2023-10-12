package service

import (
	"context"
	"errors"
	"geekgo-webook/internal/domain"
	"geekgo-webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	Profile(ctx context.Context, uid int64) (domain.User, error)
	Login(ctx context.Context, email, password string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error)
}

//type userServiceV1 struct {
//	repo repository.UserRepository
//	logger *zap.Logger // 这里指定了zap.Logger 最好是提供一个logger的接口 再提供zap的实现，以便之后扩展
//}
//
//func NewUserServiceV1(repo repository.UserRepository, logger *zap.Logger) UserService{
//	return &userServiceV1{
//		repo, repo,
//		logger: logger
//}
//}

//type userServiceV2 struct {
//	repo repository.UserRepository
//	logger logger.LoggerV1 // 这里指定了zap.Logger 最好是提供一个logger的接口 再提供zap的实现，以便之后扩展
//}
//
//func NewUserServiceV2(repo repository.UserRepository, logger  logger.LoggerV1 ) UserService{
//	return &userServiceV2{
//		repo, repo,
//		logger: logger
//}
//}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, info.OpenID)
	return u, err
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	// 数据库中保存的应该是加密的密码 在service层对密码进行加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	return u, err
}

func (svc *userService) Profile(ctx context.Context, uid int64) (domain.User, error) {
	u, err := svc.repo.FindById(ctx, uid)
	return u, err
}

func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	// 需要从数据库中找到用户 （mysql 或 redis ) 通过repo进行接口封装
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 比较密码了
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// DEBUG
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil

}
