package ioc

import (
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
	"goworkwebook/webook003/config"
	"goworkwebook/webook003/internal/repository/dao"
	"goworkwebook/webook003/pkg/gormx"
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

	// 使用prometheus.New函数创建一个新的prometheus实例，并传入配置参数
	err = db.Use(prometheus.New(prometheus.Config{
		// 数据库名称
		DBName: "webook",
		// 刷新间隔时间
		RefreshInterval: 15,
		// 指标收集器
		MetricsCollector: []prometheus.MetricsCollector{
			// MySQL指标收集器
			&prometheus.MySQL{
				// 需要收集的变量名称
				VariableNames: []string{"thread_running"},
			},
		},
	}))
	// 如果出现错误，则抛出panic
	if err != nil {
		panic(err)
	}

	// 创建一个新的回调函数，用于统计 GORM 的数据库查询
	cb := gormx.NewCallbacks(prometheus2.SummaryOpts{
		// 设置命名空间
		Namespace: "geektime_daming",
		// 设置子系统
		Subsystem: "webook",
		// 设置指标名称
		Name: "gorm_db",
		// 设置帮助信息
		Help: "统计 GORM 的数据库查询",
		// 设置常量标签
		ConstLabels: map[string]string{
			"instance_id": "my_instance",
		},
		// 设置目标值
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})
	// 使用回调函数
	err = db.Use(cb)
	// 如果出现错误，则抛出异常
	if err != nil {
		panic(err)
	}

	// 初始化数据库的表
	err = dao.InitTable(db)
	// 如果初始化失败，则抛出异常
	if err != nil {
		panic(err)
	}
	return db
}
