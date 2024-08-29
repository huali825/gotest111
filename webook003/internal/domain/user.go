package domain

import "time"

type DMUser struct {
	Id       int64
	Email    string
	Password string
	Ctime    time.Time
}
