package model

const AccountTableName = "accounts"

// Account 账户，类似支付宝，银行卡
type Account struct {
	Model
	UserID       uint    `json:"user_id" form:"user_id"`             // 创建人
	Name         string  `json:"name" form:"name"`                   // 账户名称
	Balance      float64 `json:"balance" form:"balance"`             // 账户余额
	CurrencyType uint    `json:"currency_type" form:"currency_type"` // 货币类型
	Type         uint    `json:"type" form:"type"`                   // 账户类型
	Icon         uint    `json:"icon" form:"icon"`                   // 账户图标
	Description  string  `json:"description" form:"description"`     // 账户描述
}

// RegisterAccountReq -
type RegisterAccountReq struct {
	UserID       uint    `json:"user_id" form:"user_id" binding:"required"`             // 创建人
	Name         string  `json:"name" form:"name" binding:"required,notblank"`          // 账户名称
	Balance      float64 `json:"balance" form:"balance"`                                // 账户余额，默认为0
	CurrencyType uint    `json:"currency_type" form:"currency_type" binding:"required"` // 货币类型
	Type         uint    `json:"type" form:"type" binding:"required"`                   // 账户类型
	Icon         uint    `json:"icon" form:"icon" binding:"required"`                   // 账户图标,不支持自定义，从已有的里面选
	Description  string  `json:"description" form:"description"`                        // 账户描述
}

// ModifyAccountReq -
type ModifyAccountReq struct {
	ID           uint     `json:"id" form:"id" binding:"required"`           //账户id
	UserID       uint     `json:"user_id" form:"user_id" binding:"required"` // 创建人
	Name         *string  `json:"name" form:"name"`                          // 账户名称
	Balance      *float64 `json:"balance" form:"balance"`                    // 账户余额
	CurrencyType *uint    `json:"currency_type" form:"currency_type"`        // 货币类型
	Type         *uint    `json:"type" form:"type"`                          // 账户类型
	Icon         *uint    `json:"icon" form:"icon"`                          // 账户图标
	Description  *string  `json:"description" form:"description"`            // 账户描述
}

// GetAccountReq -
type GetAccountReq struct {
	ID           *uint   `json:"id" form:"id"`                              //账户id
	UserID       uint    `json:"user_id" form:"user_id" binding:"required"` // 创建人
	Name         *string `json:"name" form:"name"`                          // 账户名称
	CurrencyType *uint   `json:"currency_type" form:"currency_type"`        // 货币类型
	Type         *uint   `json:"type" form:"type"`                          // 账户类型
	Icon         *uint   `json:"icon" form:"icon"`                          // 账户图标
}

// DeleteAccountReq -
type DeleteAccountReq struct {
	ID     uint `json:"id" form:"id" binding:"required"`           //账户id
	UserID uint `json:"user_id" form:"user_id" binding:"required"` // 创建人
}
