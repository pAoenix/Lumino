package service

import (
	"Lumino/model"
	"Lumino/store"
)

// TransactionService -
type TransactionService struct {
	TransactionStore *store.TransactionStore
	AccountBookStore *store.AccountBookStore
}

// NewTransactionService -
func NewTransactionService(transactionStore *store.TransactionStore, accountBookStore *store.AccountBookStore) *TransactionService {
	return &TransactionService{
		TransactionStore: transactionStore,
		AccountBookStore: accountBookStore,
	}
}

// Register -
func (s *TransactionService) Register(transaction *model.Transaction) error {
	return s.TransactionStore.Register(transaction)
}

// Get -
func (s *TransactionService) Get(transactionReq *model.TransactionReq) (resp []model.Transaction, err error) {
	return s.TransactionStore.Get(transactionReq)
}

// Modify -
func (s *TransactionService) Modify(transaction *model.Transaction) error {
	return s.TransactionStore.Modify(transaction)
}

// Delete -
func (s *TransactionService) Delete(transaction *model.Transaction) error {
	return s.TransactionStore.Delete(transaction)
}
