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
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
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
