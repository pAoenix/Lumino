package service

import (
	"Lumino/model"
	"Lumino/store"
	"time"
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
func (s *TransactionService) Register(transaction *model.RegisterTransactionReq) error {
	if transaction.Date.IsZero() {
		transaction.Date = time.Now()
	}
	return s.TransactionStore.Register(transaction)
}

// Get -
func (s *TransactionService) Get(transactionReq *model.GetTransactionReq) (resp []model.Transaction, err error) {
	return s.TransactionStore.Get(transactionReq)
}

// Modify -
func (s *TransactionService) Modify(transaction *model.ModifyTransactionReq) error {
	return s.TransactionStore.Modify(transaction)
}

// Delete -
func (s *TransactionService) Delete(transaction *model.DeleteTransactionReq) error {
	return s.TransactionStore.Delete(transaction)
}
