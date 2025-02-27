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

// Modify -
func (s *AccountBookService) Modify(accountBook *model.AccountBook) error {
	return s.AccountBookStore.Modify(accountBook)
}

// Get -
func (s *AccountBookService) Get(accountBookReq *model.AccountBookReq) (resp []model.AccountBook, err error) {
	return s.AccountBookStore.Get(accountBookReq)
}

// Delete -
func (s *AccountBookService) Delete(accountBook *model.AccountBook) error {
	return s.AccountBookStore.Delete(accountBook)
}

// Merge -
func (s *AccountBookService) Merge(mergeAccountBookReq *model.MergeAccountBookReq) error {
	return s.AccountBookStore.Merge(mergeAccountBookReq)
}
