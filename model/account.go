package model

import (
	"gorm.io/gorm"
)

// Account 账户，类似支付宝，银行卡
type Account struct {
	gorm.Model
	UserID      uint // 创建人
	Name        string
	Balance     float64
	Type        int
	Icon        string
	Description string
}
