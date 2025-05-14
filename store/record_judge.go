package store

import (
	"Lumino/common/http_error_code"
	"Lumino/model"
	"errors"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// ParamsJudge 判断输入参数的数据是否存在，不允许无中生有的东西存在
func ParamsJudge(db *DB, AccountBookID *uint, userIDs *pq.Int32Array, userID *uint, categoryID *uint, transactionID *uint, accountID *uint) error {
	if AccountBookID != nil {
		accountBook := model.AccountBook{}
		if err := db.Model(&model.AccountBook{}).Where("id = ?", *AccountBookID).First(&accountBook).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return http_error_code.BadRequest("账本不存在")
			}
			return http_error_code.Internal("服务内部错误",
				http_error_code.WithInternal(err))
		}
	}
	if userIDs != nil {
		var count int64
		if err := db.Model(&model.User{}).Where("id IN ?", []int32(*userIDs)).Count(&count).Error; err != nil {
			return http_error_code.Internal("服务内部错误",
				http_error_code.WithInternal(err))
		}
		if count != int64(len(*userIDs)) {
			return http_error_code.BadRequest("请确保用户都存在")
		}
	}
	if userID != nil {
		user := model.User{}
		if err := db.Model(&model.User{}).Where("id = ?", *userID).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return http_error_code.BadRequest("用户不存在")
			}
			return http_error_code.Internal("服务内部错误",
				http_error_code.WithInternal(err))
		}
	}
	if categoryID != nil {
		category := model.Category{}
		if err := db.Model(&model.Category{}).Where("id = ?", *categoryID).First(&category).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return http_error_code.BadRequest("图标不存在")
			}
			return http_error_code.Internal("服务内部错误",
				http_error_code.WithInternal(err))
		}
	}
	if transactionID != nil {
		transaction := model.Transaction{}
		if err := db.Model(&model.Transaction{}).Where("id = ?", *transactionID).First(&transaction).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return http_error_code.BadRequest("交易记录不存在")
			}
			return http_error_code.Internal("服务内部错误",
				http_error_code.WithInternal(err))
		}
	}
	if accountID != nil {
		account := model.Account{}
		if err := db.Model(&model.Account{}).Where("id = ?", *accountID).First(&account).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return http_error_code.BadRequest("账户不存在")
			}
			return http_error_code.Internal("服务内部错误",
				http_error_code.WithInternal(err))
		}
	}
	return nil
}
