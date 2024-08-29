package repository

import (
	"context"
	"goworkwebook/webook003/internal/domain"
	"goworkwebook/webook003/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	//ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.DMUser) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}
