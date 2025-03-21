package model

import (
	"gorm.io/gorm"
	"time"
)

const IncomeType = 0
const SpendingType = 1

// Transaction 交易记录
type Transaction struct {
	gorm.Model
	Type           int       // 类型:收入/支出
	Amount         float64   // 交易数额
	Date           time.Time // 日期
	UserID         uint      // 记录账户id
	CategoryID     uint      // 关联消费场景分类ID
	Description    string    // 注释
	AccountBookID  uint      // 对应的账本id
	RelatedUserIDs []uint    `gorm:"type:json"` // 涉及那些人
}

type TransactionReq struct {
	UserID        uint       // 账户id
	AccountBookID uint       // 对应的账本id
	BeginTime     *time.Time // 起始时间
	EndTime       *time.Time // 结束时间
}
