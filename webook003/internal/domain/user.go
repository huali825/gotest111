package domain

import "time"

// DMUser 在代码中使用的结构体
type DMUser struct {
	Id       int64
	Email    string
	Password string

	Nickname string
	// YYYY-MM-DD
	Birthday time.Time
	AboutMe  string

	Phone string

	// UTC 0 的时区
	Ctime time.Time

	WechatInfo WechatInfo

	//Addr Address
}
