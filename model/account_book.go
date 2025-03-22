package model

import "github.com/lib/pq"

// AccountBook -
type AccountBook struct {
	Model
	CreatorID int           `json:"creator_id" form:"creator_id"`                 // 创建人
	UserID    pq.Int32Array `gorm:"type:integer[]" json:"user_id" form:"user_id"` // 账单用户列表
	Name      string        `json:"name" form:"name"`                             // 账本名称
	Spending  float64       `json:"spending" form:"spending"`                     // 账本花费
	Income    float64       `json:"income" form:"income"`                         // 账本收入
}

// AccountBookReq -
type AccountBookReq struct {
	UserID   uint `json:"user_id" form:"user_id"`     // 用户
	ID       uint `json:"id" form:"id"`               // 账本id
	SortType int  `json:"sort_type" form:"sort_type"` // 排序模式  0: 创建时间升序，1:创建时间降序
}

// MergeAccountBookReq -
type MergeAccountBookReq struct {
	Model
	MergeAccountBookID  uint `json:"merge_account_book_id" form:"merge_account_book_id"`   // 合并的账本id  A
	MergedAccountBookID uint `json:"merged_account_book_id" form:"merged_account_book_id"` // 被合并的账本id B -> A，B的记录全部合入到A
}

// AccountBookResp -
type AccountBookResp struct {
	AccountBooks         []AccountBook `json:"account_books" form:"account_books"`                     //账本列表
	Users                []User        `json:"users" form:"users"`                                     // 涉及的用户信息
	DefaultAccountBookID uint          `json:"default_account_book_id" form:"default_account_book_id"` // 默认版本id
}
