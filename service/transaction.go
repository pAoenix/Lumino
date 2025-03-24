package service

import (
	"Lumino/common"
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
	if !common.ContainsInt(common.ConvertArrayToIntSlice(transaction.RelatedUserIDs), int(transaction.CreatorID)) {
		transaction.RelatedUserIDs = append(transaction.RelatedUserIDs, int32(transaction.CreatorID))
	}
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
