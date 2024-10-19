package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"goworkwebook/webook003/config"
	"goworkwebook/webook003/internal/repository/dao"
)

func InitDB() *gorm.DB {

	//db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:30003)/webook"))
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))

	if err != nil {
		// 我只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
