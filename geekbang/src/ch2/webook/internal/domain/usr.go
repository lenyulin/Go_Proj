package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Password string
	//UTC 0
	Ctime        time.Time
	Utime        time.Time
	NickName     string
	Birthday     string
	Introduction string
}
