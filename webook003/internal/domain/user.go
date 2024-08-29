package domain

import "time"

// DMUser 在代码中使用的结构体
type DMUser struct {
	Id       int64
	Email    string
	Password string
	Ctime    time.Time
}
