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
	if err := tx.Model(&model.Transaction{}).Create(transaction).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	// 更新账本数值
	if transaction.Type == model.IncomeType {
		if err := tx.Model(&model.AccountBook{}).
			Where("id = ?", accountBook.ID).
			Update("income", accountBook.Income+transaction.Amount).Error; err != nil {
			tx.Rollback() // 回滚事务
			return err
		}
	} else {
		if err := tx.Model(&model.AccountBook{}).Update("spending", accountBook.Spending+transaction.Amount).Error; err != nil {
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
	tx := s.db.Begin()
	accountBook := model.AccountBook{}
	// 加锁
	if err := tx.Model(model.AccountBook{}).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", transaction.AccountBookID).
		First(&accountBook).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	// 更新账本数值
	tmpTransaction := model.Transaction{}
	if err := tx.Model(&model.Transaction{}).
		Select("*").Where("id = ?", transaction.ID).First(&tmpTransaction).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 更新后的交易信息入账
	if transaction.Type == model.IncomeType {
		// 如果修改前后都是收入
		if tmpTransaction.Type == model.IncomeType {
			if err := tx.Model(&model.AccountBook{}).
				Where("id = ?", accountBook.ID).
				Update("income", accountBook.Income+transaction.Amount-tmpTransaction.Amount).Error; err != nil {
				tx.Rollback() // 回滚事务
				return err
			}
		} else { // 如果支出变收入
			if err := tx.Model(&model.AccountBook{}).
				Where("id = ?", accountBook.ID).
				Update("spending", accountBook.Income-tmpTransaction.Amount).
				Update("income", accountBook.Income+transaction.Amount).Error; err != nil {
				tx.Rollback() // 回滚事务
				return err
			}
		}
	} else {
		// 收入变支出
		if tmpTransaction.Type == model.IncomeType {
			if err := tx.Model(&model.AccountBook{}).
				Where("id = ?", accountBook.ID).
				Update("income", accountBook.Income-tmpTransaction.Amount).
				Update("spending", accountBook.Spending+transaction.Amount).Error; err != nil {
				tx.Rollback() // 回滚事务
				return err
			}
		} else {
			// 一直是支出
			if err := tx.Model(&model.AccountBook{}).
				Where("id = ?", accountBook.ID).
				Update("spending", accountBook.Spending+transaction.Amount-tmpTransaction.Amount).Error; err != nil {
				tx.Rollback() // 回滚事务
				return err
			}
		}
	}
	// 交易信息更新
	if err := tx.Model(&model.Transaction{}).Where("id = ?", transaction.ID).Updates(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// Delete -
func (s *TransactionStore) Delete(transaction *model.Transaction) error {
	tx := s.db.Begin()
	accountBook := model.AccountBook{}
	// 加锁
	if err := tx.Model(model.AccountBook{}).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", transaction.AccountBookID).
		First(&accountBook).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	// 更新账本数值
	tmpTransaction := model.Transaction{}
	if err := tx.Model(&model.Transaction{}).
		Select("*").Where("id = ?", transaction.ID).First(&tmpTransaction).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 把修改前的信息，从账本退去
	if tmpTransaction.Type == model.IncomeType {
		if err := tx.Model(&model.AccountBook{}).
			Where("id = ?", accountBook.ID).
			Update("income", accountBook.Income-tmpTransaction.Amount).Error; err != nil {
			tx.Rollback() // 回滚事务
			return err
		}
	} else {
		if err := tx.Model(&model.AccountBook{}).
			Where("id = ?", accountBook.ID).
			Update("spending", accountBook.Spending-tmpTransaction.Amount).Error; err != nil {
			tx.Rollback() // 回滚事务
			return err
		}
	}

	if err := tx.Model(&model.Transaction{}).Where("id = ?", transaction.ID).Delete(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
