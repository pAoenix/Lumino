package model

import "github.com/lib/pq"

const AccountBookTableName = "account_books"

// AccountBook -
type AccountBook struct {
	Model
	CreatorID uint          `json:"creator_id" form:"creator_id"`                                               // 创建人,不允许修改
	UserIDs   pq.Int32Array `gorm:"type:integer[]" json:"user_ids" form:"user_ids" swaggertype:"array,integer"` // 账本用户列表
	Name      string        `json:"name" form:"name"`                                                           // 账本名称
	Spending  float64       `json:"spending" form:"spending"`                                                   // 账本花费
	Income    float64       `json:"income" form:"income"`                                                       // 账本收入
}

// RegisterAccountBookReq -
type RegisterAccountBookReq struct {
	Name      string        `json:"name" form:"name" binding:"required,notblank"`                               // 账本名称
	UserIDs   pq.Int32Array `gorm:"type:integer[]" json:"user_ids" form:"user_ids" swaggertype:"array,integer"` // 账本用户列表
	CreatorID uint          `json:"creator_id" form:"creator_id" binding:"required"`                            // 创建人
}

// GetAccountBookReq -
type GetAccountBookReq struct {
	CreatorID uint  `json:"user_id" form:"user_id" binding:"required"` // 用户
	ID        *uint `json:"id" form:"id" swaggerignore:"true"`         // 账本id
	SortType  int   `json:"sort_type" form:"sort_type"`                // 排序模式  0: 创建时间升序，1:创建时间降序
}

// MergeAccountBookReq -
type MergeAccountBookReq struct {
	MergeAccountBookID  uint `json:"merge_account_book_id" form:"merge_account_book_id" binding:"required"`   // 合并的账本id  A
	MergedAccountBookID uint `json:"merged_account_book_id" form:"merged_account_book_id" binding:"required"` // 被合并的账本id B -> A，B的记录全部合入到A
	CreatorID           uint `json:"creator_id" form:"creator_id" binding:"required"`                         // 创建人(只能合并自己创建的账本)
}

// AccountBookResp -
type AccountBookResp struct {
	AccountBooks         []AccountBook `json:"account_books" form:"account_books"`                     //账本列表
	Users                []User        `json:"users" form:"users"`                                     // 涉及的用户信息
	DefaultAccountBookID uint          `json:"default_account_book_id" form:"default_account_book_id"` // 默认版本id
}

// ModifyAccountBookReq -
type ModifyAccountBookReq struct {
	ID      uint           `json:"id" form:"id" binding:"required"`                                            // 账本id
	Name    string         `json:"name" form:"name" binding:"omitempty,notblank"`                              // 账本名称
	UserIDs *pq.Int32Array `gorm:"type:integer[]" json:"user_ids" form:"user_ids" swaggertype:"array,integer"` // 账本用户列表
}

// DeleteAccountBookReq -
type DeleteAccountBookReq struct {
	ID uint `json:"id" form:"id" binding:"required"` // 账本id
}
