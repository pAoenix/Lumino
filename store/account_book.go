package store

import (
	"Lumino/common"
	"Lumino/model"
	"errors"
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
func (s *AccountBookStore) Register(accountBook *model.RegisterAccountBookReq) error {
	// 判断用户是否存在
	if err := s.db.Model(&model.User{}).Where("id = ?", accountBook.CreatorID).First(&model.User{}).Error; err != nil {
		return errors.New("用户不存在" + err.Error())
	}
	var count int64
	if err := s.db.Model(&model.User{}).Where("id IN ?", []int32(accountBook.UserIDs)).Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(accountBook.UserIDs)) {
		return errors.New("请确保所有用户均存在")
	}
	newAB := model.AccountBook{
		Name:      accountBook.Name,
		CreatorID: accountBook.CreatorID,
		UserIDs:   accountBook.UserIDs,
	}
	return s.db.Model(&model.AccountBook{}).Create(&newAB).Error
}

// Get -
func (s *AccountBookStore) Get(accountBookReq *model.GetAccountBookReq) (resp []model.AccountBook, err error) {
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
func (s *AccountBookStore) Modify(accountBookReq *model.ModifyAccountBookReq) error {
	// 判断用户是否存在
	if accountBookReq.UserIDs != nil {
		var count int64
		if err := s.db.Model(&model.User{}).Where("id IN ?", []int32(*accountBookReq.UserIDs)).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(*accountBookReq.UserIDs)) {
			return errors.New("请确保所有用户均存在")
		}
		// 需要保证创建人在用户里
		accountBook := model.AccountBook{}
		if err := s.db.Model(&model.AccountBook{}).Where("id = ?", accountBookReq.ID).First(&accountBook).Error; err != nil {
			return err
		}
		if !common.ContainsInt(common.ConvertArrayToIntSlice(*accountBookReq.UserIDs), int(accountBook.CreatorID)) {
			*accountBookReq.UserIDs = append(*accountBookReq.UserIDs, int32(accountBook.CreatorID))
		}
	}
	return s.db.Model(&model.AccountBook{}).Where("id = ?", accountBookReq.ID).Updates(&accountBookReq).Error
}

// Delete -
func (s *AccountBookStore) Delete(accountBookReq *model.DeleteAccountBookReq) error {
	return s.db.Model(&model.AccountBook{}).Delete(&model.AccountBook{Model: model.Model{ID: accountBookReq.ID}}).Error
}

// Merge -
func (s *AccountBookStore) Merge(mergeAccountBookReq *model.MergeAccountBookReq) error {
	tx := s.db.Begin()
	// 账本加锁
	accountBook := model.AccountBook{}
	if err := tx.Model(&model.AccountBook{}).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", mergeAccountBookReq.MergeAccountBookID).
		First(&accountBook).Error; err != nil {
		tx.Rollback() // 回滚事务
		return errors.New("合并账本不存在" + err.Error())
	}
	// 更新账本数值
	var transactions []model.Transaction
	if err := tx.Model(&model.Transaction{}).Select("*").Where("account_book_id = ?", mergeAccountBookReq.MergedAccountBookID).Find(&transactions).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	// 获取被合并的账本的用户，需要对用户合并
	mergedAB := model.AccountBook{}
	if err := tx.Model(&model.AccountBook{}).
		Where("id = ?", mergeAccountBookReq.MergedAccountBookID).
		First(&mergedAB).Error; err != nil {
		tx.Rollback() // 回滚事务
		return errors.New("被合并账本不存在" + err.Error())
	}
	for _, userID := range mergedAB.UserIDs {
		if !common.ContainsInt(common.ConvertArrayToIntSlice(accountBook.UserIDs), int(userID)) {
			accountBook.UserIDs = append(accountBook.UserIDs, userID)
		}
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
		Updates(model.AccountBook{
			Spending: accountBook.Spending + spending,
			Income:   accountBook.Income + income,
			UserIDs:  accountBook.UserIDs,
		}).Error; err != nil {
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
	if err := tx.Model(&model.AccountBook{}).
		Delete(&model.AccountBook{Model: model.Model{ID: mergeAccountBookReq.MergedAccountBookID}}).
		Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
