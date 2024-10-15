package service

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"goworkwebook/webook003/internal/domain"
	"goworkwebook/webook003/internal/repository"
)

var ErrDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

type UserService interface {
	Signup(ctx context.Context, u domain.DMUser) error
	Login(ctx context.Context, email string, password string) (domain.DMUser, error)
	UpdateNonSensitiveInfo(ctx context.Context,
		user domain.DMUser) error
	FindById(ctx context.Context,
		uid int64) (domain.DMUser, error)
	FindOrCreate(ctx context.Context, phone string) (domain.DMUser, error)
	FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.DMUser, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) Signup(ctx context.Context, u domain.DMUser) error {
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

func (svc *userService) Login(ctx context.Context, email string, password string) (domain.DMUser, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	// 根据传入的email查找用户
	if errors.Is(err, repository.ErrUserNotFound) {
		// 如果找不到用户，返回空用户和错误信息
		return domain.DMUser{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		// 如果发生其他错误，返回空用户和错误信息
		return domain.DMUser{}, err //系统错误 等会儿再来谈论
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	// 将用户密码和传入的密码进行比对
	if err != nil {
		// 如果密码不匹配，返回错误信息
		return domain.DMUser{}, ErrInvalidUserOrPassword
	}
	return u, err
}

// FindOrCreate 先找，找不到就创建
func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.DMUser, error) {
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

func (svc *userService) UpdateNonSensitiveInfo(ctx context.Context,
	user domain.DMUser) error {
	// UpdateNicknameAndXXAnd
	return svc.repo.UpdateNonZeroFields(ctx, user)
}

func (svc *userService) FindById(ctx context.Context,
	uid int64) (domain.DMUser, error) {
	return svc.repo.FindById(ctx, uid)
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.DMUser, error) {
	u, err := svc.repo.FindByWechat(ctx, wechatInfo.OpenId)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	// 这边就是意味着是一个新用户
	// JSON 格式的 wechatInfo
	zap.L().Info("新用户", zap.Any("wechatInfo", wechatInfo))
	//svc.logger.Info("新用户", zap.Any("wechatInfo", wechatInfo))
	err = svc.repo.Create(ctx, domain.DMUser{
		WechatInfo: wechatInfo,
	})
	if err != nil && !errors.Is(err, repository.ErrUserDuplicateEmail) {
		return domain.DMUser{}, err
	}
	return svc.repo.FindByWechat(ctx, wechatInfo.OpenId)
}
