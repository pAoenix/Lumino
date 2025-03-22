package store

import (
	"Lumino/model"
	"gorm.io/gorm/clause"
)

// AccountBookStore -
type AccountBookStore struct {
	db *DB
}

// NewAccountBookStore -
func NewAccountBookStore(db *DB) *AccountBookStore {
	return &AccountBookStore{
		db: db,
	}
}

// Register -
func (s *AccountBookStore) Register(accountBook *model.AccountBook) error {
	return s.db.Model(&model.AccountBook{}).Create(accountBook).Error
}

// Get -
func (s *AccountBookStore) Get(accountBookReq *model.AccountBookReq) (resp []model.AccountBook, err error) {
	sort := "id desc"
	if accountBookReq.SortType == 0 {
		sort = "id"
	}
	if s.db.Model(&model.AccountBook{}).Order(sort).Where("? = any(user_ids)", accountBookReq.UserID).Find(&resp).Error != nil {
		return nil, err
	} else {
		return
	}
}

// Modify -
func (s *AccountBookStore) Modify(accountBookReq *model.AccountBook) error {
	return s.db.Model(&model.AccountBook{}).Where("id = ?", accountBookReq.ID).Updates(&accountBookReq).Error
}

// Delete -
func (s *AccountBookStore) Delete(accountBookReq *model.AccountBook) error {
	return s.db.Model(&model.AccountBook{}).Where(accountBookReq).Delete(&accountBookReq).Error
}

// Merge -
func (s *AccountBookStore) Merge(mergeAccountBookReq *model.MergeAccountBookReq) error {
	tx := s.db.Begin()
	// 账本加锁
	accountBook := model.AccountBook{}
	if err := tx.Model(model.AccountBook{}).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", mergeAccountBookReq.MergeAccountBookID).
		First(&accountBook).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	// 更新账本数值
	var transactions []model.Transaction
	if err := tx.Model(&model.Transaction{}).Select("*").Where("account_book_id = ?", mergeAccountBookReq.MergedAccountBookID).Find(&transactions).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	income := 0.
	spending := 0.
	for _, transaction := range transactions {
		if transaction.Type == model.IncomeType {
			income += transaction.Amount
		} else {
			spending += transaction.Amount
		}
	}
	if err := tx.Model(&model.AccountBook{}).Where("id = ?", accountBook.ID).
		Select("spending", "income").
		Updates(model.AccountBook{Spending: accountBook.Spending + spending, Income: accountBook.Income + income}).
		Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	// 合并账本交易明细
	if err := tx.Model(&model.Transaction{}).
		Where("account_book_id = ?", mergeAccountBookReq.MergedAccountBookID).
		Update("account_book_id", mergeAccountBookReq.MergeAccountBookID).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 删除账本
	if err := tx.Model(model.AccountBook{}).
		Delete(&model.AccountBook{Model: model.Model{ID: mergeAccountBookReq.MergedAccountBookID}}).
		Error; err != nil {
		return err
	}
	return tx.Commit().Error
}
