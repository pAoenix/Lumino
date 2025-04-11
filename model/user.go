package model

import "github.com/lib/pq"

const UserTableName = "users"

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

// GetUserReq -
type GetUserReq struct {
	ID uint `json:"id" form:"id" binding:"required"` // 用户id
}

// DeleteUserReq -
type DeleteUserReq struct {
	ID uint `json:"id" form:"id" binding:"required"` // 用户id
}

// ModifyUserReq -
type ModifyUserReq struct {
	ID                   uint               `json:"id" form:"id" binding:"required"`
	Name                 string             `json:"name" form:"name" gorm:"uniqueIndex:idx_user_name"`                          //账号名称，昵称，全局唯一
	DefaultAccountBookID uint               `json:"default_account_book_id" form:"default_account_book_id"`                     // 默认账本id
	Balance              float64            `json:"balance" form:"balance"`                                                     // 余额
	Friend               pq.Int32Array      `json:"friend" form:"friend" gorm:"type:integer[]" swaggertype:"array,integer"`     // 朋友列表
	IconUrl              string             `json:"icon_url" form:"icon_url" swaggerignore:"true"`                              // 用户头像的对象存储地址
	BalanceDetail        map[string]float64 `json:"balance_detail" form:"balance_detail" gorm:"type:json" swaggerignore:"true"` // 余额详情
}

// RegisterUserReq -
type RegisterUserReq struct {
	Name string `json:"name" form:"name" binding:"required"` //账号名称，昵称，全局唯一
}
