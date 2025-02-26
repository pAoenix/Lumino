package service

import (
	"Lumino/model"
	"Lumino/store"
)

// TransactionService -
type TransactionService struct {
	TransactionStore *store.TransactionStore
}

// NewTransactionService -
func NewTransactionService(TransactionStore *store.TransactionStore) *TransactionService {
	return &TransactionService{
		TransactionStore: TransactionStore,
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
