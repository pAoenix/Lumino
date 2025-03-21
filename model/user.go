package model

// User 用户账户
type User struct {
	Model
	Name                 string             `form:"name" form:"name"`                                       //账号名称，昵称
	DefaultAccountBookID uint               `json:"default_account_book_id" form:"default_account_book_id"` // 默认账本id
	Balance              float64            `json:"balance" form:"balance"`                                 // 余额
	BalanceDetail        map[string]float64 `json:"balance_detail" form:"balance_detail" gorm:"type:json"`  // 余额详情
	Friend               []int              `json:"friend" form:"friend" gorm:"type:json"`                  // 朋友列表
	Icon                 string             `json:"icon" form:"icon"`                                       // 用户头像的对象存储地址
}
