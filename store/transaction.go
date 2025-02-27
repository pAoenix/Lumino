package store

import "Lumino/model"

// TransactionStore -
type TransactionStore struct {
	db *DB
}

// NewTransactionStore -
func NewTransactionStore(db *DB) *TransactionStore {
	return &TransactionStore{
		db: db,
	}
}

// Register -
func (s *TransactionStore) Register(transaction *model.Transaction) error {
	return s.db.Create(transaction).Error
}

// Get -
func (s *TransactionStore) Get(transactionReq *model.TransactionReq) (resp []model.Transaction, err error) {
	if s.db.Where(transactionReq).Find(&resp).Error != nil {
		return nil, err
	} else {
		return
	}
}

// Modify -
func (s *TransactionStore) Modify(transaction *model.Transaction) error {
	return s.db.Updates(transaction).Error
}

// Delete -
func (s *TransactionStore) Delete(transaction *model.Transaction) error {
	return s.db.Delete(transaction).Error
}
