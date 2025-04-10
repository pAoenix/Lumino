package model

// Account 账户，类似支付宝，银行卡
type Account struct {
	Model
	UserID      uint    `json:"user_id" form:"user_id"`         // 创建人
	Name        string  `json:"name" form:"name"`               // 账户名称
	Balance     float64 `json:"balance" form:"balance"`         // 账户余额
	Type        int     `json:"type" form:"type"`               // 账户类型
	Icon        string  `json:"icon" form:"icon"`               // 账户图标
	Description string  `json:"description" form:"description"` // 账户描述
}
