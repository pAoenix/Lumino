package model

import "gorm.io/gorm"

// AccountBook -
type AccountBook struct {
	gorm.Model
	CreatorID int    // 创建人
	UserId    []int  // 账单用户列表
	Name      string // 账本名称
}

// AccountBookReq -
type AccountBookReq struct {
	UserId   int // 用户
	SortType int // 排序模式  0: 创建时间升序，1:创建时间降序
}

// MergeAccountBookReq -
type MergeAccountBookReq struct {
	gorm.Model
	MergeAccountBookID  int // 合并的账本id  A
	MergedAccountBookID int // 被合并的账本id B -> A，B的记录全部合入到A
}
