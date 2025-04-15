package model

import (
	"github.com/lib/pq"
	"time"
)

const TransactionTableName = "transactions"
const IncomeType = 1
const SpendingType = 2 //

// Transaction 交易记录
type Transaction struct {
	Model
	Type          int       `json:"type" form:"type"`                       // 类型:收入/支出
	Amount        float64   `json:"amount" form:"amount"`                   // 交易数额
	Date          time.Time `json:"date" form:"date"`                       // 日期
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
}

// GetTransactionReq -
type GetTransactionReq struct {
	ID            *uint      `json:"id" form:"id"`                                              // 交易id
	AccountBookID uint       `json:"account_book_id" form:"account_book_id" binding:"required"` // 对应的账本id
	BeginTime     *time.Time `json:"begin_time" form:"begin_time"`                              // 起始时间
	EndTime       *time.Time `json:"end_time" form:"end_time"`                                  // 结束时间
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
