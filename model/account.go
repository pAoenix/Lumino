package model

// Account 账户，类似支付宝，银行卡
type Account struct {
	Model
	UserID      uint    `json:"user_id" form:"user_id"` // 创建人
	Name        string  `json:"name" form:"name"`
	Balance     float64 `json:"balance" form:"balance"`
	Type        int     `json:"type" form:"type"`
	Icon        string  `json:"icon" form:"icon"`
	Description string  `json:"description" form:"description"`
}
