package service

import (
	"Lumino/model"
	"Lumino/store"
)

// AccountService -
type AccountService struct {
	AccountStore *store.AccountStore
}

// NewAccountService -
func NewAccountService(AccountStore *store.AccountStore) *AccountService {
	return &AccountService{
		AccountStore: AccountStore,
	}
}

// Register -
func (s *AccountService) Register(Account *model.RegisterAccountReq) (account model.Account, err error) {
	return s.AccountStore.Register(Account)
}

// Modify -
func (s *AccountService) Modify(Account *model.ModifyAccountReq) (account model.Account, err error) {
	return s.AccountStore.Modify(Account)
}

// Get -
func (s *AccountService) Get(Account *model.GetAccountReq) (account []model.Account, err error) {
	return s.AccountStore.Get(Account)
}

// Delete -
func (s *AccountService) Delete(Account *model.DeleteAccountReq) error {
	return s.AccountStore.Delete(Account)
}
