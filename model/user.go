package model

import (
	"gorm.io/gorm"
)

// User 用户账户
type User struct {
	gorm.Model
	Name                 string             //账号名称，昵称
	DefaultAccountBookID int                // 默认账本id
	Balance              float64            // 余额
	BalanceDetail        map[string]float64 `gorm:"type:json"` // 余额详情
	Friend               []int              // 朋友列表
	IconUrl              string             // 用户头像的对象存储地址
}
