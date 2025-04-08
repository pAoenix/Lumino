package model

import "github.com/lib/pq"

// User 用户账户
type User struct {
	Model
	Name                 string  `json:"name" form:"name" gorm:"uniqueIndex:idx_user_name"`      //账号名称，昵称，全局唯一
	DefaultAccountBookID uint    `json:"default_account_book_id" form:"default_account_book_id"` // 默认账本id
	Balance              float64 `json:"balance" form:"balance"`                                 // 余额
	// BalanceDetail 余额详情
	// @swagger:type object
	// @additionalProperties type=number format=double
	// @example {"temperature":36.5,"humidity":0.42}
	BalanceDetail map[string]float64 `json:"balance_detail" form:"balance_detail" gorm:"type:json"`
	Friend        pq.Int32Array      `json:"friend" form:"friend" gorm:"type:integer[]" swaggertype:"array,integer"` // 朋友列表
	IconUrl       string             `json:"icon_url" form:"icon_url"`                                               // 用户头像的对象存储地址
}

// GetUser -
type GetUser struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

// DeleteUser -
type DeleteUser struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

// ModifyUser -
type ModifyUser struct {
	ID                   uint    `json:"id" form:"id" binding:"required"`
	Name                 string  `json:"name" form:"name" gorm:"uniqueIndex:idx_user_name"`      //账号名称，昵称，全局唯一
	DefaultAccountBookID uint    `json:"default_account_book_id" form:"default_account_book_id"` // 默认账本id
	Balance              float64 `json:"balance" form:"balance"`                                 // 余额
	// BalanceDetail 余额详情
	// @swagger:type object
	// @additionalProperties type=number format=double
	// @example {"temperature":36.5,"humidity":0.42}
	BalanceDetail map[string]float64 `json:"balance_detail" form:"balance_detail" gorm:"type:json"`
	Friend        pq.Int32Array      `json:"friend" form:"friend" gorm:"type:integer[]" swaggertype:"array,integer"` // 朋友列表
	IconUrl       string             `json:"icon_url" form:"icon_url"`                                               // 用户头像的对象存储地址
}
