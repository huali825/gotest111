package dao

import "gorm.io/gorm"

// 在数据库中初始化table
func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}