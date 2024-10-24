package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art IsDaoArticle) (int64, error)
	UpdateById(ctx context.Context, entity IsDaoArticle) error
	Sync(ctx context.Context, entity IsDaoArticle) (int64, error)
	SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]IsDaoArticle, error)
}

type ArticleGORMDAO struct {
	db *gorm.DB
}

func NewArticleGORMDAO(db *gorm.DB) ArticleDAO {
	return &ArticleGORMDAO{
		db: db,
	}
}

func (a *ArticleGORMDAO) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]IsDaoArticle, error) {
	var arts []IsDaoArticle
	err := a.db.WithContext(ctx).
		Where("author_id = ?", uid).
		Offset(offset).Limit(limit).
		// a ASC, B DESC
		Order("utime DESC").
		Find(&arts).Error
	return arts, err
}

func (a *ArticleGORMDAO) SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error {
	now := time.Now().UnixMilli()
	return a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&IsDaoArticle{}).
			Where("id = ? and author_id = ?", uid, id).
			Updates(map[string]any{
				"utime":  now,
				"status": status,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return errors.New("ID 不对或者创作者不对")
		}
		return tx.Model(&PublishedArticle{}).
			Where("id = ?", uid).
			Updates(map[string]any{
				"utime":  now,
				"status": status,
			}).Error
	})
}

// Sync 函数用于同步文章数据，根据文章的id判断是更新还是插入
func (a *ArticleGORMDAO) Sync(ctx context.Context, art IsDaoArticle) (int64, error) {
	// 获取文章的id
	var id = art.Id
	// 开始事务
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var (
			err error
		)
		// 创建一个新的ArticleGORMDAO实例
		dao := NewArticleGORMDAO(tx)
		// 如果id大于0，则更新文章
		if id > 0 {
			err = dao.UpdateById(ctx, art)
			// 否则，插入文章
		} else {
			id, err = dao.Insert(ctx, art)
		}
		// 如果有错误，则返回错误
		if err != nil {
			return err
		}
		// 设置文章的id
		art.Id = id
		// 获取当前时间
		now := time.Now().UnixMilli()
		// 创建一个PublishedArticle实例
		pubArt := PublishedArticle(art)
		// 设置创建时间和更新时间
		pubArt.Ctime = now
		pubArt.Utime = now
		// 在事务中创建PublishedArticle实例，如果id冲突，则更新
		err = tx.Clauses(clause.OnConflict{
			// 对MySQL不起效，但是可以兼容别的方言
			// INSERT xxx ON DUPLICATE KEY SET `title`=?
			// 别的方言：
			// sqlite INSERT XXX ON CONFLICT DO UPDATES WHERE
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":   pubArt.Title,
				"content": pubArt.Content,
				"utime":   now,
				"status":  pubArt.Status,
			}),
		}).Create(&pubArt).Error
		// 返回错误
		return err
	})
	// 返回文章的id和错误
	return id, err
}
func (a *ArticleGORMDAO) SyncV1(ctx context.Context, art IsDaoArticle) (int64, error) {
	tx := a.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	// 防止后面业务panic
	defer tx.Rollback()

	var (
		id  = art.Id
		err error
	)
	dao := NewArticleGORMDAO(tx)
	if id > 0 {
		err = dao.UpdateById(ctx, art)
	} else {
		id, err = dao.Insert(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	now := time.Now().UnixMilli()
	pubArt := PublishedArticle(art)
	pubArt.Ctime = now
	pubArt.Utime = now
	err = tx.Clauses(clause.OnConflict{
		// 对MySQL不起效，但是可以兼容别的方言
		// INSERT xxx ON DUPLICATE KEY SET `title`=?
		// 别的方言：
		// sqlite INSERT XXX ON CONFLICT DO UPDATES WHERE
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   pubArt.Title,
			"content": pubArt.Content,
			"utime":   now,
		}),
	}).Create(&pubArt).Error
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, nil
}

func (a *ArticleGORMDAO) UpdateById(ctx context.Context, art IsDaoArticle) error {
	now := time.Now().UnixMilli()
	res := a.db.WithContext(ctx).Model(&art).
		Where("id = ? AND author_id = ?", art.Id, art.AuthorId).Updates(map[string]any{
		"title":   art.Title,
		"content": art.Content,
		"status":  art.Status,
		"utime":   now,
	})
	if res.Error != nil {
		return res.Error
	}
	// 我怎么知道有没有更新数据？
	if res.RowsAffected == 0 {
		// 创作者不对，说明有人在瞎搞
		return errors.New("ID 不对或者创作者不对")
	}
	return nil
}

// Insert 在ArticleGORMDAO结构体中定义Insert方法，用于向数据库中插入一条Article记录
func (a *ArticleGORMDAO) Insert(ctx context.Context, art IsDaoArticle) (int64, error) {
	// 获取当前时间戳
	now := time.Now().UnixMilli()
	// 将当前时间戳赋值给Article的创建时间和更新时间
	art.Ctime = now
	art.Utime = now
	// 使用WithContext方法将上下文传递给数据库，并使用Create方法创建一条Article记录
	err := a.db.WithContext(ctx).Create(&art).Error
	// 返回Article的Id和错误信息
	return art.Id, err
}

//type IsDaoArticle struct {
//	Id      int64  `gorm:"primaryKey,autoIncrement"`
//	Title   string `gorm:"type=varchar(4096)"`
//	Content string `gorm:"type=BLOB"`
//	// 我要根据创作者ID来查询
//	AuthorId int64 `gorm:"index"`
//
//	Status uint8
//
//	Ctime int64
//	// 更新时间
//	Utime int64
//}

type IsDaoArticle struct {
	Id      int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title   string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	Content string `gorm:"type=BLOB" bson:"content,omitempty"`
	// 我要根据创作者ID来查询
	AuthorId int64 `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8 `bson:"status,omitempty"`
	Ctime    int64 `bson:"ctime,omitempty"`
	// 更新时间
	Utime int64 `bson:"utime,omitempty"`
}

// PublishedArticle 是一个结构体，用于表示已发布的文章
// 这里的两个写法表示如何cp一个已有的结构体
type PublishedArticle IsDaoArticle

type PublishedArticleV1 struct {
	IsDaoArticle
}
