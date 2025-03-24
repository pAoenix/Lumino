package store

import "Lumino/model"

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
func (s *AccountStore) Register(account *model.Account) error {
	return s.db.Model(model.Account{}).Create(account).Error
}

// Modify -
func (s *AccountStore) Modify(account *model.Account) error {
	return s.db.Model(model.Account{}).Where(account).Updates(account).Error
}

// Get -
func (s *AccountStore) Get(accountReq *model.Account) (account []model.Account, err error) {
	if err = s.db.Model(model.Account{}).Where(accountReq).Find(&account).Error; err != nil {
		return []model.Account{}, err
	}
	return
}

// Delete -
func (s *AccountStore) Delete(account *model.Account) error {
	return s.db.Model(model.Account{}).Delete(account).Error
}
