package store

import (
	"Lumino/model"
	"gorm.io/gorm/clause"
)

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
	tx := s.db.Begin()
	accountBook := model.AccountBook{}
	// 加锁
	if err := tx.Model(model.AccountBook{}).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", transaction.AccountBookID).
		First(&accountBook).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	// 新建交易记录
	if err := s.db.Model(&model.Transaction{}).Create(transaction).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	// 更新账本数值
	if transaction.Type == model.IncomeType {
		if err := s.db.Model(&model.AccountBook{}).Update("income", accountBook.Income+transaction.Amount).Error; err != nil {
			tx.Rollback() // 回滚事务
			return err
		}
	} else {
		if err := s.db.Model(&model.AccountBook{}).Update("spending", accountBook.Spending+transaction.Amount).Error; err != nil {
			tx.Rollback() // 回滚事务
			return err
		}
	}
	return tx.Commit().Error
}

// Get -
func (s *TransactionStore) Get(transactionReq *model.TransactionReq) (resp []model.Transaction, err error) {
	sql := s.db.Model(&model.Transaction{})
	if transactionReq.BeginTime != nil {
		sql.Where("date >= ", &transactionReq.BeginTime)
	}
	if transactionReq.EndTime != nil {
		sql.Where("date <= ", &transactionReq.EndTime)
	}
	if sql.Where(transactionReq).Find(&resp).Error != nil {
		return nil, err
	} else {
		return
	}
}

// Modify -
func (s *TransactionStore) Modify(transaction *model.Transaction) error {
	return s.db.Model(&model.Transaction{}).Updates(transaction).Error
}

// Delete -
func (s *TransactionStore) Delete(transaction *model.Transaction) error {
	return s.db.Model(&model.Transaction{}).Delete(transaction).Error
}
