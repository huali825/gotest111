package dao

import (
	"context"
	"gorm.io/gorm"
)

type ArticleAuthorDAO interface {
	Create(ctx context.Context, art IsDaoArticle) (int64, error)
	Update(ctx context.Context, art IsDaoArticle) error
}

type ArticleGORMAuthorDAO struct {
	db *gorm.DB
}

func NewArticleGORMAuthorDAO(db *gorm.DB) ArticleAuthorDAO {
	return &ArticleGORMAuthorDAO{
		db: db,
	}
}

func (a *ArticleGORMAuthorDAO) Create(ctx context.Context, art IsDaoArticle) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (a *ArticleGORMAuthorDAO) Update(ctx context.Context, art IsDaoArticle) error {
	//TODO implement me
	panic("implement me")
}
