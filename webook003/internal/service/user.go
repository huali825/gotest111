package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
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

func (svc *UserService) Login(c *gin.Context, email string, password string) error {
	u, err := svc.repo.FindByMail(c, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return ErrInvalidUserOrPassword
	}
	if err != nil {
		return err //系统错误 等会儿再来谈论
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// 密码错误
		return ErrInvalidUserOrPassword
	}
	return err
}
