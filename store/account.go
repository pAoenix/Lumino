package store

import (
	"Lumino/common/http_error_code"
	"Lumino/model"
	"github.com/jinzhu/copier"
)

// AccountStore -
type AccountStore struct {
	db *DB
}

// NewAccountStore -
func NewAccountStore(db *DB) *AccountStore {
	return &AccountStore{
		db: db,
	}
}

// Register -
func (s *AccountStore) Register(accountReq *model.RegisterAccountReq) (account model.Account, err error) {
	if err = ParamsJudge(s.db, nil, nil,
		&accountReq.UserID, nil, nil, nil); err != nil {
		return
	}
	if err = copier.Copy(&account, &accountReq); err != nil {
		return account, http_error_code.Internal("服务内异常",
			http_error_code.WithInternal(err))
	}
	err = s.db.Model(&model.Account{}).Create(&account).Error
	return
}

// Modify -
func (s *AccountStore) Modify(accountReq *model.ModifyAccountReq) (account model.Account, err error) {
	if err = ParamsJudge(s.db, nil, nil,
		&accountReq.UserID, nil, nil, &accountReq.ID); err != nil {
		return
	}
	if err = copier.Copy(&account, &accountReq); err != nil {
		return account, http_error_code.Internal("服务内异常",
			http_error_code.WithInternal(err))
	}
	// 需要判断账户和用户一致
	tmpAccount := model.Account{}
	if err = s.db.Model(&model.Account{}).
		Where("id = ? and user_id = ?", accountReq.ID, accountReq.UserID).
		First(&tmpAccount).Error; err != nil {
		return account, http_error_code.BadRequest("账户不存在本人名下",
			http_error_code.WithInternal(err))
	}
	err = s.db.Model(&model.Account{}).Where("id = ?", accountReq.ID).Updates(&account).Find(&account).Error
	return
}

// Get -
func (s *AccountStore) Get(accountReq *model.GetAccountReq) (account []model.Account, err error) {
	if err = ParamsJudge(s.db, nil, nil,
		&accountReq.UserID, nil, nil, accountReq.ID); err != nil {
		return
	}
	if accountReq.ID != nil {
		tmpAccount := model.Account{}
		if err = s.db.Model(&model.Account{}).
			Where("id = ? and user_id = ?", accountReq.ID, accountReq.UserID).
			First(&tmpAccount).Error; err != nil {
			return account, http_error_code.BadRequest("账户不存在本人名下",
				http_error_code.WithInternal(err))
		}
	}
	err = s.db.Model(&model.Account{}).Where(&accountReq).Find(&account).Error
	return
}

// Delete -
func (s *AccountStore) Delete(accountReq *model.DeleteAccountReq) error {
	if err := ParamsJudge(s.db, nil, nil,
		&accountReq.UserID, nil, nil, &accountReq.ID); err != nil {
		return err
	}
	tmpAccount := model.Account{}
	if err := s.db.Model(&model.Account{}).
		Where("id = ? and user_id = ?", accountReq.ID, accountReq.UserID).
		First(&tmpAccount).Error; err != nil {
		return http_error_code.BadRequest("账户不存在本人名下",
			http_error_code.WithInternal(err))
	}
	return s.db.Model(&model.Account{}).Where("id = ?", accountReq.ID).Delete(&model.Account{}).Error
}
