package dao

import "gorm.io/gorm"

// InitTable 在数据库中初始化table
func InitTable(db *gorm.DB) error {
	//return db.AutoMigrate(&User{}, &IsDaoArticle{}, &PublishedArticle{})
	return db.AutoMigrate(
		&User{},
		&IsDaoArticle{},
		&PublishedArticle{},
		&Interactive{},
		&UserLikeBiz{})
}
