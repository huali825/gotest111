package gormx

import (
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"time"
)

// Callbacks 结构体，用于记录数据库操作的时间
type Callbacks struct {
	vector *prometheus.SummaryVec
}

// NewCallbacks 方法，创建一个新的 Callbacks 实例
func NewCallbacks(opts prometheus.SummaryOpts) *Callbacks {
	vector := prometheus.NewSummaryVec(opts,
		[]string{"type", "table"})
	prometheus.MustRegister(vector)
	return &Callbacks{
		vector: vector,
	}
}

// Name 方法，返回回调的名称
func (c *Callbacks) Name() string {
	return "prometheus"
}

// Initialize 方法，初始化回调函数
func (c *Callbacks) Initialize(db *gorm.DB) error {
	// 在创建操作之前注册回调函数
	err := db.Callback().Create().Before("*").
		Register("prometheus_create_before", c.Before())
	if err != nil {
		return err
	}

	// 在创建操作之后注册回调函数
	err = db.Callback().Create().After("*").
		Register("prometheus_create_after", c.After("CREATE"))
	if err != nil {
		return err
	}

	// 在查询操作之前注册回调函数
	err = db.Callback().Query().Before("*").
		Register("prometheus_query_before", c.Before())
	if err != nil {
		return err
	}

	// 在查询操作之后注册回调函数
	err = db.Callback().Query().After("*").
		Register("prometheus_query_after", c.After("QUERY"))
	if err != nil {
		return err
	}

	// 在原始查询操作之前注册回调函数
	err = db.Callback().Query().Before("*").
		Register("prometheus_raw_before", c.Before())
	if err != nil {
		return err
	}

	// 在原始查询操作之后注册回调函数
	err = db.Callback().Raw().After("*").
		Register("prometheus_raw_after", c.After("RAW"))
	if err != nil {
		return err
	}

	// 在更新操作之前注册回调函数
	err = db.Callback().Update().Before("*").
		Register("prometheus_update_before", c.Before())
	if err != nil {
		return err
	}

	// 在更新操作之后注册回调函数
	err = db.Callback().Update().After("*").
		Register("prometheus_update_after", c.After("UPDATE"))
	if err != nil {
		return err
	}

	// 在删除操作之前注册回调函数
	err = db.Callback().Delete().Before("*").
		Register("prometheus_delete_before", c.Before())
	if err != nil {
		return err
	}

	// 在删除操作之后注册回调函数
	err = db.Callback().Update().After("*").
		Register("prometheus_delete_after", c.After("DELETE"))
	if err != nil {
		return err
	}

	// 在行操作之前注册回调函数
	err = db.Callback().Row().Before("*").
		Register("prometheus_row_before", c.Before())
	if err != nil {
		return err
	}

	// 在行操作之后注册回调函数
	err = db.Callback().Update().After("*").
		Register("prometheus_row_after", c.After("ROW"))
	return err
}

// Before 方法，记录操作开始的时间
func (c *Callbacks) Before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		start := time.Now()
		db.Set("start_time", start)
	}
}

// After 方法，记录操作结束的时间，并计算操作耗时
func (c *Callbacks) After(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		val, _ := db.Get("start_time")
		start, ok := val.(time.Time)
		if ok {
			duration := time.Since(start).Milliseconds()
			c.vector.WithLabelValues(typ, db.Statement.Table).
				Observe(float64(duration))
		}
	}
}
