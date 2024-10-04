package repository

import (
	"context"
	"database/sql"
	"goworkwebook/webook003/internal/domain"
	"goworkwebook/webook003/internal/repository/cache"
	"goworkwebook/webook003/internal/repository/dao"
	"time"
)

var (
	// ErrUserDuplicateEmail errors.New("邮箱冲突")
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail

	// ErrUserNotFound 没找到
	ErrUserNotFound = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.DMUser) error {
	return r.dao.Insert(ctx, r.toEntity(u))
}

func (r *UserRepository) FindByMail(ctx context.Context, email string) (domain.DMUser, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.DMUser{}, err
	}
	return r.toDomain(u), nil
}

func (r *UserRepository) toDomain(u dao.User) domain.DMUser {
	return domain.DMUser{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
	}
}

// toEntity 将domain转换为dao
func (r *UserRepository) toEntity(u domain.DMUser) dao.User {
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

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.DMUser, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.DMUser{}, err
	}
	return r.toDomain(u), nil
}
