package repository

import (
	"context"
	"github.com/gin-gonic/gin"
	"goworkwebook/webook003/internal/domain"
	"goworkwebook/webook003/internal/repository/cache"
	"goworkwebook/webook003/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail

	//errors.New("未找到用户")
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
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindByMail(c *gin.Context, email string) (domain.DMUser, error) {
	u, err := r.dao.FindByEmail(c, email)
	if err != nil {
		return domain.DMUser{}, err
	}
	return r.toDomain(u), nil
}

func (r *UserRepository) toDomain(u dao.User) domain.DMUser {
	return domain.DMUser{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}
}
