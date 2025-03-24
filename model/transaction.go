package model

import (
	"github.com/lib/pq"
	"time"
)

const IncomeType = 0
const SpendingType = 1

// Transaction 交易记录
type Transaction struct {
	Model
	Type           int           `json:"type" form:"type"`                                               // 类型:收入/支出
	Amount         float64       `json:"amount" form:"amount"`                                           // 交易数额
	Date           time.Time     `json:"date" form:"date"`                                               // 日期
	CreatorID      uint          `json:"creator_id" form:"creator_id"`                                   // 记录账户id
	CategoryID     uint          `json:"category_id" form:"category_id"`                                 // 关联消费场景分类ID
	Description    string        `json:"description" form:"description"`                                 // 注释
	AccountBookID  uint          `json:"account_book_id" form:"account_book_id"`                         // 对应的账本id
	RelatedUserIDs pq.Int32Array `json:"related_user_ids" form:"related_user_ids" gorm:"type:integer[]"` // 涉及那些人
}

type TransactionReq struct {
	UserID        uint       `json:"user_id" form:"user_id"`                 // 账户id
	AccountBookID uint       `json:"account_book_id" form:"account_book_id"` // 对应的账本id
	BeginTime     *time.Time `json:"begin_time" form:"begin_time"`           // 起始时间
	EndTime       *time.Time `json:"end_time" form:"end_time"`               // 结束时间
}
