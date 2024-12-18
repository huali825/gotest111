package service

import (
	"context"
	"errors"
	"fmt"
	"goworkwebook/webook003/internal/repository"
	"goworkwebook/webook003/internal/service/sms"
	"math/rand"
)

var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context,
		biz, phone, inputCode string) (bool, error)
}

type codeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &codeService{
		repo: repo,
		sms:  smsSvc,
	}
}

func (svc *codeService) Send(ctx context.Context, biz, phone string) error {
	code := svc.generate()
	err := svc.repo.Set(ctx, biz, phone, code)
	// 你在这儿，是不是要开始发送验证码了？
	if err != nil {
		return err
	}
	const codeTplId = "1877556"
	return svc.sms.Send(ctx, codeTplId, []string{code}, phone)
}

// Verify函数用于验证输入的验证码是否正确
func (svc *codeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	// 调用repo的Verify函数进行验证
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	// 如果验证次数过多，则返回false和nil
	if errors.Is(err, repository.ErrCodeVerifyTooMany) {
		// 相当于，我们对外面屏蔽了验证次数过多的错误，我们就是告诉调用者，你这个不对
		return false, nil
	}
	// 否则返回验证结果和错误信息
	return ok, err
}

func (svc *codeService) generate() string {
	// 0-999999
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
