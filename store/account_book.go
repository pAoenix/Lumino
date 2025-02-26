package store

import "Lumino/model"

// AccountBookStore -
type AccountBookStore struct {
	db *DB
}

// NewAccountBookStore -
func NewAccountBookStore(db *DB) *AccountBookStore {
	return &AccountBookStore{
		db: db,
	}
}

// Register -
func (s *AccountBookStore) Register(accountBook *model.AccountBook) error {
	return s.db.Create(accountBook).Error
}

// Get -
func (s *AccountBookStore) Get(accountBookReq *model.AccountBookReq) (resp []model.AccountBook, err error) {
	if s.db.Where(accountBookReq).Find(&resp).Error != nil {
		return nil, err
	} else {
		return
	}
}
