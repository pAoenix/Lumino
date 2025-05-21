package model

import (
	"github.com/lib/pq"
	"time"
)

const TransactionTableName = "transactions"
const IncomeType = 1
const SpendingType = 2 //
const Epsilon = 1e-6

// Transaction 交易记录
type Transaction struct {
	Model
	Type          int       `json:"type" form:"type"`                       // 类型:收入/支出
	Amount        float64   `json:"amount" form:"amount"`                   // 交易数额
	Date          time.Time `json:"date" form:"date" gorm:"index"`          // 日期
	CreatorID     uint      `json:"creator_id" form:"creator_id"`           // 创建人id
	PayUserID     uint      `json:"pay_user_id" form:"pay_user_id"`         // 付款人id
	CategoryID    uint      `json:"category_id" form:"category_id"`         // 关联消费场景分类ID
	Description   string    `json:"description" form:"description"`         // 注释
	AccountBookID uint      `json:"account_book_id" form:"account_book_id"` // 对应的账本id
	// RelatedUserIDs 涉及那些人
	RelatedUserIDs pq.Int32Array `json:"related_user_ids" form:"related_user_ids" gorm:"type:integer[]" swaggertype:"array,integer"`
}

// RegisterTransactionReq -
type RegisterTransactionReq struct {
	Type          int       `json:"type" form:"type" binding:"required,oneof=1 2" `            // 类型:收入/支出
	Amount        float64   `json:"amount" form:"amount" binding:"required,gt=0"`              // 交易数额
	Date          time.Time `json:"date" form:"date"`                                          // 日期
	CreatorID     *uint     `json:"creator_id" form:"creator_id" binding:"required"`           // 创建人id
	PayUserID     *uint     `json:"pay_user_id" form:"pay_user_id" binding:"required"`         // 付款人id
	CategoryID    *uint     `json:"category_id" form:"category_id" binding:"required"`         // 关联消费场景分类ID
	Description   string    `json:"description" form:"description"`                            // 注释
	AccountBookID *uint     `json:"account_book_id" form:"account_book_id" binding:"required"` // 对应的账本id
	//  RelatedUserIDs 涉及那些人
	RelatedUserIDs *pq.Int32Array `json:"related_user_ids" form:"related_user_ids" binding:"required,min=1" gorm:"type:integer[]" swaggertype:"array,integer"`
	AccountID      *uint          `json:"account_id" form:"account_id"` // 对应的账户，可不填
}

// GetTransactionReq -
type GetTransactionReq struct {
	ID            *uint      `json:"id" form:"id"`                                              // 交易id
	UserID        *uint      `json:"user_id" form:"user_id" binding:"required"`                 // 用户id
	AccountBookID *uint      `json:"account_book_id" form:"account_book_id" binding:"required"` // 对应的账本id
	BeginTime     *time.Time `json:"begin_time" form:"begin_time"`                              // 起始时间
	EndTime       *time.Time `json:"end_time" form:"end_time"`                                  // 结束时间
	Type          *int       `json:"type" form:"type" binding:"omitempty,oneof=1 2"`            // 类型:收入/支出
	CategoryID    *uint      `json:"category_id" form:"category_id"`                            // 类别
}

// DailyTransaction 按天分组的消费数据
type DailyTransaction struct {
	Date     string        // 日期，格式如 "2023-03-15"
	Items    []Transaction // 当天的记账条目
	Spending float64       // 当天总支持
	Income   float64       // 当天总收入
}

// TransactionResp -
type TransactionResp struct {
	Transactions []DailyTransaction `json:"transactions" form:"transactions"` //账本列表
	Users        []User             `json:"users" form:"users"`               // 涉及的用户信息
	Categorys    []Category         `json:"categorys" form:"categorys"`       // 图标信息
}

// DeleteTransactionReq -
type DeleteTransactionReq struct {
	ID uint `json:"id" form:"id" binding:"required"` //交易记录id
}

// ModifyTransactionReq -
type ModifyTransactionReq struct {
	ID            uint      `json:"id" form:"id" binding:"required"`                                     //交易记录id
	AccountBookID *uint     `json:"account_book_id" form:"account_book_id"`                              // 对应的账本id
	Type          int       `json:"type" form:"type" binding:"required_with=Amount,omitempty,oneof=1 2"` // 类型:收入/支出,交易类型必须是0或者1
	Amount        float64   `json:"amount" form:"amount" binding:"required_with=Type,omitempty,gt=0"`    // 交易数额,数值需要>0
	Date          time.Time `json:"date" form:"date"`                                                    //交易日期
	CategoryID    *uint     `json:"category_id" form:"category_id"`                                      // 关联消费场景分类ID
	PayUserID     *uint     `json:"pay_user_id" form:"pay_user_id"`                                      // 付款人id
	Description   string    `json:"description" form:"description"`                                      // 注释
	// RelatedUserIDs 涉及那些人
	RelatedUserIDs *pq.Int32Array `json:"related_user_ids" form:"related_user_ids" gorm:"type:integer[]" swaggertype:"array,integer"`
}
