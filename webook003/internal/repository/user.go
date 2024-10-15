package repository

import (
	"context"
	"database/sql"
	"goworkwebook/webook003/internal/domain"
	"goworkwebook/webook003/internal/repository/cache"
	"goworkwebook/webook003/internal/repository/dao"
	"log"
	"time"
)

var (
	// ErrUserDuplicateEmail errors.New("邮箱冲突")
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail

	// ErrUserNotFound 没找到
	ErrUserNotFound = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.DMUser) error
	FindByEmail(ctx context.Context, email string) (domain.DMUser, error)
	UpdateNonZeroFields(ctx context.Context, user domain.DMUser) error
	FindByPhone(ctx context.Context, phone string) (domain.DMUser, error)
	FindById(ctx context.Context, uid int64) (domain.DMUser, error)
	FindByWechat(ctx context.Context, openId string) (domain.DMUser, error)
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *CachedUserRepository) Create(ctx context.Context, u domain.DMUser) error {
	return r.dao.Insert(ctx, r.toEntity(u))
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.DMUser, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.DMUser{}, err
	}
	return r.toDomain(u), nil
}

func (r *CachedUserRepository) toDomain(u dao.User) domain.DMUser {
	return domain.DMUser{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
		Ctime:    time.UnixMilli(u.Ctime),
	}
}

// toEntity 将domain转换为dao
func (r *CachedUserRepository) toEntity(u domain.DMUser) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}

func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.DMUser, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.DMUser{}, err
	}
	return r.toDomain(u), nil
}

func (repo *CachedUserRepository) UpdateNonZeroFields(ctx context.Context,
	user domain.DMUser) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(user))
}

func (repo *CachedUserRepository) FindById(ctx context.Context, uid int64) (domain.DMUser, error) {
	du, err := repo.cache.Get(ctx, uid)
	// 只要 err 为 nil，就返回
	if err == nil {
		return du, nil
	}

	// err 不为 nil，就要查询数据库
	// err 有两种可能
	// 1. key 不存在，说明 redis 是正常的
	// 2. 访问 redis 有问题。可能是网络有问题，也可能是 redis 本身就崩溃了

	u, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.DMUser{}, err
	}
	du = repo.toDomain(u)

	err = repo.cache.Set(ctx, du)
	if err != nil {
		// 网络崩了，也可能是 redis 崩了
		log.Println(err)
	}
	return du, nil
}

func (repo *CachedUserRepository) FindByWechat(ctx context.Context, openId string) (domain.DMUser, error) {
	ue, err := repo.dao.FindByWechat(ctx, openId)
	if err != nil {
		return domain.DMUser{}, err
	}
	return repo.toDomain(ue), nil
}
