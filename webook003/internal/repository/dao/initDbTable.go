package dao

import "gorm.io/gorm"

// InitTable 在数据库中初始化table
func InitTable(db *gorm.DB) error {
	//return db.AutoMigrate(&User{}, &IsDaoArticle{}, &PublishedArticle{})
	return db.AutoMigrate(
		&User{},              //用户表
		&IsDaoArticle{},      //草稿文章表
		&PublishedArticle{},  //已发布文章表
		&Interactive{},       //阅读 点赞 收藏 计数表
		&UserLikeBiz{},       //点赞表
		&UserCollectionBiz{}, //收藏表
	)

}
