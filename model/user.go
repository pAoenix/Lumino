package model

import "gorm.io/gorm"

// User 用户账户
type User struct {
	gorm.Model
	Name string //账号名称，昵称
}
