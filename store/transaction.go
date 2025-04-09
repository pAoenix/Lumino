package store

import (
	"Lumino/model"
	"errors"
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
func (s *TransactionStore) Register(transactionSeq *model.RegisterTransactionReq) error {
	tx := s.db.Begin()
	accountBook := model.AccountBook{}
	// 判断用户是否存在
	if err := tx.Model(&model.User{}).Where("id = ?", transactionSeq.CreatorID).First(&model.User{}).Error; err != nil {
		tx.Rollback() // 回滚事务
		return errors.New("用户不存在" + err.Error())
	}
	// 判断图标是否存在
	if err := tx.Model(&model.Category{}).Where("id = ?", transactionSeq.CategoryID).First(&model.Category{}).Error; err != nil {
		tx.Rollback() // 回滚事务
		return errors.New("图标不存在" + err.Error())
	}
	// 加锁
	if err := tx.Model(model.AccountBook{}).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", transactionSeq.AccountBookID).
		First(&accountBook).Error; err != nil {
		tx.Rollback() // 回滚事务
		return errors.New("账本不存在" + err.Error())
	}
	// 新建交易记录
	transaction := model.Transaction{
		Date:           transactionSeq.Date,
		Type:           transactionSeq.Type,
		Amount:         transactionSeq.Amount,
		CreatorID:      transactionSeq.CreatorID,
		CategoryID:     transactionSeq.CategoryID,
		Description:    transactionSeq.Description,
		AccountBookID:  transactionSeq.AccountBookID,
		RelatedUserIDs: transactionSeq.RelatedUserIDs,
	}
	if err := tx.Model(&model.Transaction{}).Create(&transaction).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	// 更新账本数值
	if transactionSeq.Type == model.IncomeType {
		if err := tx.Model(&model.AccountBook{}).
			Where("id = ?", accountBook.ID).
			Update("income", accountBook.Income+transactionSeq.Amount).Error; err != nil {
			tx.Rollback() // 回滚事务
			return err
		}
	} else {
		if err := tx.Model(&model.AccountBook{}).
			Where("id = ?", accountBook.ID).
			Update("spending", accountBook.Spending+transactionSeq.Amount).Error; err != nil {
			tx.Rollback() // 回滚事务
			return err
		}
	}
	return tx.Commit().Error
}

// Get -
func (s *TransactionStore) Get(transactionReq *model.GetTransactionReq) (resp []model.Transaction, err error) {
	sql := s.db.Model(&model.Transaction{})
	if transactionReq.BeginTime != nil {
		sql.Where("date >= ?", &transactionReq.BeginTime)
	}
	if transactionReq.EndTime != nil {
		sql.Where("date <= ?", &transactionReq.EndTime)
	}
	if transactionReq.UserID != 0 {
		sql.Where("? = ANY(related_user_ids)", &transactionReq.UserID)
	}
	if sql.Where("account_book_id = ?", transactionReq.AccountBookID).Find(&resp).Error != nil {
		return nil, err
	} else {
		return
	}
}

// Modify -
func (s *TransactionStore) Modify(modifyTransaction *model.ModifyTransactionReq) error {
	tx := s.db.Begin()
	accountBook := model.AccountBook{}
	transaction := model.Transaction{}
	// 判断图标是否存在
	if modifyTransaction.CategoryID != 0 {
		if err := tx.Model(&model.Category{}).Where("id = ?", modifyTransaction.CategoryID).First(&model.Category{}).Error; err != nil {
			tx.Rollback() // 回滚事务
			return errors.New("图标不存在" + err.Error())
		}
	}
	if modifyTransaction.AccountBookID != 0 {
		if err := tx.Model(&model.AccountBook{}).Where("id = ?", modifyTransaction.AccountBookID).First(&model.AccountBook{}).Error; err != nil {
			tx.Rollback() // 回滚事务
			return errors.New("更新账本不存在" + err.Error())
		}
	}
	// 交易记录查询
	if err := tx.Model(model.Transaction{}).Where("id = ?", modifyTransaction.ID).
		First(&transaction).Error; err != nil {
		tx.Rollback()
		return errors.New("交易记录不存在" + err.Error())
	}
	// 账本加锁
	if err := tx.Model(model.AccountBook{}).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", transaction.AccountBookID).
		First(&accountBook).Error; err != nil {
		tx.Rollback() // 回滚事务
		return errors.New("账本不存在" + err.Error())
	}
	// TODO 用户判断
	if modifyTransaction.Amount > 0 {
		// 重新查询，避免脏读
		if err := tx.Model(model.Transaction{}).Where("id = ?", modifyTransaction.ID).
			First(&transaction).Error; err != nil {
			tx.Rollback()
			return errors.New("交易记录不存在" + err.Error())
		}
		// 更新后的交易信息入账
		if modifyTransaction.Type == model.IncomeType {
			// 如果修改前后都是收入
			if transaction.Type == model.IncomeType {
				if err := tx.Model(&model.AccountBook{}).
					Where("id = ?", accountBook.ID).
					Update("income", accountBook.Income+modifyTransaction.Amount-transaction.Amount).Error; err != nil {
					tx.Rollback() // 回滚事务
					return err
				}
			} else { // 如果支出变收入
				if err := tx.Model(&model.AccountBook{}).
					Where("id = ?", accountBook.ID).
					Update("spending", accountBook.Spending-transaction.Amount).
					Update("income", accountBook.Income+modifyTransaction.Amount).Error; err != nil {
					tx.Rollback() // 回滚事务
					return err
				}
			}
		} else {
			// 收入变支出
			if transaction.Type == model.IncomeType {
				if err := tx.Model(&model.AccountBook{}).
					Where("id = ?", accountBook.ID).
					Update("income", accountBook.Income-transaction.Amount).
					Update("spending", accountBook.Spending+modifyTransaction.Amount).Error; err != nil {
					tx.Rollback() // 回滚事务
					return err
				}
			} else {
				// 一直是支出
				if err := tx.Model(&model.AccountBook{}).
					Where("id = ?", accountBook.ID).
					Update("spending", accountBook.Spending+modifyTransaction.Amount-transaction.Amount).Error; err != nil {
					tx.Rollback() // 回滚事务
					return err
				}
			}
		}
	}
	// 交易信息更新
	if err := tx.Model(&model.Transaction{}).Where("id = ?", modifyTransaction.ID).Updates(modifyTransaction).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// Delete -
func (s *TransactionStore) Delete(transaction *model.DeleteTransactionReq) error {
	tx := s.db.Begin()
	accountBook := model.AccountBook{}
	// 加锁
	if err := tx.Model(model.AccountBook{}).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", transaction.AccountBookID).
		First(&accountBook).Error; err != nil {
		tx.Rollback() // 回滚事务
		return errors.New("账本不存在" + err.Error())
	}
	// 更新账本数值
	tmpTransaction := model.Transaction{}
	if err := tx.Model(&model.Transaction{}).
		Where(transaction).First(&tmpTransaction).Error; err != nil {
		tx.Rollback()
		return errors.New("记录信息有误，账本和交易记录不匹配" + err.Error())
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

	if err := tx.Model(&model.Transaction{}).Where(transaction).Delete(&model.Transaction{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
