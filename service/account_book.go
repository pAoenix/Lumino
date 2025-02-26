package service

import (
	"Lumino/model"
	"Lumino/store"
)

// AccountBookService -
type AccountBookService struct {
	AccountBookStore *store.AccountBookStore
}

// NewAccountBookService -
func NewAccountBookService(AccountBookStore *store.AccountBookStore) *AccountBookService {
	return &AccountBookService{
		AccountBookStore: AccountBookStore,
	}
}

// Register -
func (s *AccountBookService) Register(accountBook *model.AccountBook) error {
	return s.AccountBookStore.Register(accountBook)
}

// Get -
func (s *AccountBookService) Get(accountBookReq *model.AccountBookReq) (resp []model.AccountBook, err error) {
	return s.AccountBookStore.Get(accountBookReq)
}
