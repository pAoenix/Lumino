package store

import (
	"Lumino/common/http_error_code"
	"Lumino/model"
	"errors"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// ParamsJudge 判断输入参数的数据是否存在，不允许无中生有的东西存在
func ParamsJudge(db *DB, AccountBookID *uint, userIDs *pq.Int32Array) error {
	if AccountBookID != nil {
		accountBook := model.AccountBook{}
		if err := db.Model(model.AccountBook{}).Where("id = ?", *AccountBookID).First(&accountBook).Error; err != nil {
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
	return nil
}
