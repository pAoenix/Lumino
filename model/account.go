package model

// Account 账户，类似支付宝，银行卡
type Account struct {
	Model
	UserID      uint    `json:"user_id"` // 创建人
	Name        string  `json:"name"`
	Balance     float64 `json:"balance"`
	Type        int     `json:"type"`
	Icon        string  `json:"icon"`
	Description string  `json:"description"`
}
