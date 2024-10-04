package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"goworkwebook/webook003/internal/domain"
	"goworkwebook/webook003/internal/repository"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail

var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.DMUser) error {
	// 你要考虑加密放在哪里的问题了
	//使用 BCrypt 加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	// 然后就是，存起来
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.DMUser, error) {
	u, err := svc.repo.FindByMail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.DMUser{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.DMUser{}, err //系统错误 等会儿再来谈论
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// 密码错误
		return domain.DMUser{}, ErrInvalidUserOrPassword
	}
	return u, err
}

// FindOrCreate 先找，找不到就创建
func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.DMUser, error) {
	// 先找一下，我们认为，大部分用户是已经存在的用户
	u, err := svc.repo.FindByPhone(ctx, phone)
	if !errors.Is(err, repository.ErrUserNotFound) {
		// 有两种情况
		// err == nil, u 是可用的
		// err != nil，系统错误，
		return u, err
	}
	// 用户没找到
	err = svc.repo.Create(ctx, domain.DMUser{
		Phone: phone,
	})
	// 有两种可能，一种是 err 恰好是唯一索引冲突（phone）
	// 一种是 err != nil，系统错误
	if err != nil && !errors.Is(err, repository.ErrUserDuplicateEmail) {
		return domain.DMUser{}, err
	}
	// 要么 err ==nil，要么ErrDuplicateUser，也代表用户存在
	// 主从延迟，理论上来讲，强制走主库
	return svc.repo.FindByPhone(ctx, phone)
}
