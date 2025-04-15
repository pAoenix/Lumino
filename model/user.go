package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/lib/pq"
)

const UserTableName = "users"

type BalanceDetail map[string]float64

// User 用户账户
type User struct {
	Model
	PhoneNumber          string        `json:"phone_number" form:"phone_number" gorm:"uniqueIndex:idx_phone_number"`   // 手机号
	Name                 string        `json:"name" form:"name" gorm:"uniqueIndex:idx_user_name"`                      //账号名称，昵称，全局唯一
	DefaultAccountBookID uint          `json:"default_account_book_id" form:"default_account_book_id"`                 // 默认账本id
	Balance              float64       `json:"balance" form:"balance"`                                                 // 余额
	Friend               pq.Int32Array `json:"friend" form:"friend" gorm:"type:integer[]" swaggertype:"array,integer"` // 朋友列表
	IconUrl              string        `json:"icon_url" form:"icon_url"`                                               // 用户头像的对象存储地址
	// BalanceDetail 余额详情
	// @swagger:type object
	// @additionalProperties type=number format=double
	// @example {"temperature":36.5,"humidity":0.42}
	BalanceDetail BalanceDetail `json:"balance_detail" form:"balance_detail" gorm:"type:json"`
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
	Name                 string             `json:"name" form:"name" gorm:"uniqueIndex:idx_user_name" binding:"omitempty,notblank"` //账号名称，昵称，全局唯一
	DefaultAccountBookID *uint              `json:"default_account_book_id" form:"default_account_book_id"`                         // 默认账本id
	Balance              float64            `json:"balance" form:"balance"`                                                         // 余额
	Friend               *pq.Int32Array     `json:"friend" form:"friend" gorm:"type:integer[]" swaggertype:"array,integer"`         // 朋友列表
	IconUrl              string             `json:"icon_url" form:"icon_url" swaggerignore:"true"`                                  // 用户头像的对象存储地址
	BalanceDetail        map[string]float64 `json:"balance_detail" form:"balance_detail" gorm:"type:json" swaggerignore:"true"`     // 余额详情
	PhoneNumber          string             `json:"phone_number" form:"phone_number" binding:"omitempty,phone"`                     // 手机号
}

// ModifyProfilePhotoReq -
type ModifyProfilePhotoReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

// RegisterUserReq -
type RegisterUserReq struct {
	Name        string `json:"name" form:"name" binding:"required,notblank"`              //账号名称，昵称，全局唯一
	PhoneNumber string `json:"phone_number" form:"phone_number" binding:"required,phone"` // 手机号
}

// Scan - 重写，支持map类型pg
func (bm *BalanceDetail) Scan(value interface{}) error {
	if value == nil {
		*bm = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("invalid balance detail data")
	}

	return json.Unmarshal(bytes, bm)
}

// Value - 重写，支持map类型pg
func (bm BalanceDetail) Value() (driver.Value, error) {
	if bm == nil {
		return nil, nil
	}
	return json.Marshal(bm)
}
