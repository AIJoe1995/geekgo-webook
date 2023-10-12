package service

import (
	"context"
	"fmt"
	"geekgo-webook/internal/repository"
	"geekgo-webook/internal/service/sms"
	"math/rand"
)

var (
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrUnknownForCode         = repository.ErrUnknownForCode
)

const codeTplId = "1877556"

// codeTplId 修改从配置文件读取 这会和viper强耦合，
// 需要考虑同时读写的并发问题， 用原子操作

//var codeTplId atomic.String = atomic.String{}

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type codeService struct {
	repo   repository.CodeRepository // 需要从repo里拿到验证码 验证 保存到repo
	smsSvc sms.Service               // 需要短信服务 发送验证码
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &codeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *codeService) Send(ctx context.Context, biz string, phone string) error {
	//codeTplId.Store("1877556")
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	codeTplId.Store(viper.GetString("code.tpl.id"))
	//})
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		// 有问题
		return err
	}
	// 这前面成功了

	// 发送出去

	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	//err = svc.smsSvc.Send(ctx, codeTplId.Load(), []string{code}, phone)
	return err
}

func (svc *codeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	// 验证验证码
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *codeService) generateCode() string {
	// 六位数，num 在 0, 999999 之间，包含 0 和 999999
	num := rand.Intn(1000000)
	// 不够六位的，加上前导 0
	// 000001
	return fmt.Sprintf("%6d", num)
}
