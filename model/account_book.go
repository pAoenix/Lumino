package model

// AccountBook -
type AccountBook struct {
	Model
	CreatorID int     // 创建人
	UserID    []int   `gorm:"type:json"` // 账单用户列表
	Name      string  // 账本名称
	Spending  float64 // 账本花费
	Income    float64 // 账本收入
}

// AccountBookReq -
type AccountBookReq struct {
	UserID   uint // 用户
	ID       uint // 账本id
	SortType int  // 排序模式  0: 创建时间升序，1:创建时间降序
}

// MergeAccountBookReq -
type MergeAccountBookReq struct {
	Model
	MergeAccountBookID  uint // 合并的账本id  A
	MergedAccountBookID uint // 被合并的账本id B -> A，B的记录全部合入到A
}

// AccountBookResp -
type AccountBookResp struct {
	AccountBooks         []AccountBook //账本列表
	Users                []User        // 涉及的用户信息
	DefaultAccountBookID uint          // 默认版本id
}
