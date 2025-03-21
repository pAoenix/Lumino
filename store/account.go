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
func (s *AccountStore) Register(Account *model.Account) error {
	return s.db.Model(model.Account{}).Create(Account).Error
}

// Modify -
func (s *AccountStore) Modify(Account *model.Account) error {
	return s.db.Model(model.Account{}).Updates(Account).Error
}

// Get -
func (s *AccountStore) Get(Account *model.Account) error {
	return s.db.Model(model.Account{}).Where(Account).Find(Account).Error
}

// Delete -
func (s *AccountStore) Delete(Account *model.Account) error {
	return s.db.Model(model.Account{}).Delete(Account).Error
}
