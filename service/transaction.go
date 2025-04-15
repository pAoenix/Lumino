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
func (s *TransactionService) Register(transactionReq *model.RegisterTransactionReq) (resp model.Transaction, err error) {
	if transactionReq.Date.IsZero() {
		transactionReq.Date = time.Now()
	}
	return s.TransactionStore.Register(transactionReq)
}

// Get -
func (s *TransactionService) Get(transactionReq *model.GetTransactionReq) (resp []model.Transaction, err error) {
	return s.TransactionStore.Get(transactionReq)
}

// Modify -
func (s *TransactionService) Modify(transactionReq *model.ModifyTransactionReq) (transaction model.Transaction, err error) {
	return s.TransactionStore.Modify(transactionReq)
}

// Delete -
func (s *TransactionService) Delete(transactionReq *model.DeleteTransactionReq) error {
	return s.TransactionStore.Delete(transactionReq)
}
