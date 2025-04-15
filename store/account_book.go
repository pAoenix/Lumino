package store

import (
	"Lumino/common"
	"Lumino/common/http_error_code"
	"Lumino/model"
	"github.com/jinzhu/copier"
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
func (s *AccountBookStore) Register(accountBookReq *model.RegisterAccountBookReq) (accountBook model.AccountBook, err error) {
	if err = ParamsJudge(s.db, nil, &accountBookReq.UserIDs, &accountBookReq.CreatorID, nil); err != nil {
		return
	}
	if err = copier.Copy(&accountBook, &accountBookReq); err != nil {
		return accountBook, http_error_code.Internal("服务内异常",
			http_error_code.WithInternal(err))
	}
	if err = s.db.Model(&model.AccountBook{}).Create(&accountBook).Error; err != nil {
		return model.AccountBook{}, http_error_code.Internal("服务内异常",
			http_error_code.WithInternal(err))
	}
	return
}

// Get -
func (s *AccountBookStore) Get(accountBookReq *model.GetAccountBookReq) (resp []model.AccountBook, err error) {
	sort := "id desc"
	if accountBookReq.SortType == 0 {
		sort = "id"
	}
	if err = ParamsJudge(s.db, &accountBookReq.ID, nil, &accountBookReq.UserID, nil); err != nil {
		return
	}

	if s.db.Model(&model.AccountBook{}).Order(sort).Where("? = any(user_ids)", accountBookReq.UserID).Find(&resp).Error != nil {
		return nil, err
	} else {
		return
	}
}

// Modify -
func (s *AccountBookStore) Modify(accountBookReq *model.ModifyAccountBookReq) (accountBook model.AccountBook, err error) {
	if err = ParamsJudge(s.db, &accountBookReq.ID, accountBookReq.UserIDs, nil, nil); err != nil {
		return
	}
	if accountBookReq.UserIDs != nil {
		// 需要保证创建人在用户里
		accountBookTmp := model.AccountBook{}
		if err = s.db.Model(&model.AccountBook{}).Where("id = ?", accountBookReq.ID).First(&accountBookTmp).Error; err != nil {
			return
		}
		if !common.ContainsInt(common.ConvertArrayToIntSlice(*accountBookReq.UserIDs), int(accountBookTmp.CreatorID)) {
			*accountBookReq.UserIDs = append(*accountBookReq.UserIDs, int32(accountBookTmp.CreatorID))
		}
	}
	if err = copier.Copy(&accountBook, &accountBookReq); err != nil {
		return accountBook, http_error_code.Internal("服务内异常",
			http_error_code.WithInternal(err))
	}
	if err = s.db.Model(&model.AccountBook{}).Where("id = ?", accountBookReq.ID).Updates(&accountBook).Find(&accountBook).Error; err != nil {
		return
	}
	return
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
		Where("id = ? and creator_id = ?", mergeAccountBookReq.MergeAccountBookID, mergeAccountBookReq.CreatorID).
		First(&accountBook).Error; err != nil {
		tx.Rollback() // 回滚事务
		return http_error_code.BadRequest("合并账本不存在本人名下",
			http_error_code.WithInternal(err))
	}
	// 获取被合并的账本的用户，需要对用户合并
	mergedAccountBook := model.AccountBook{}
	if err := tx.Model(&model.AccountBook{}).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? and creator_id = ?", mergeAccountBookReq.MergedAccountBookID, mergeAccountBookReq.CreatorID).
		First(&mergedAccountBook).Error; err != nil {
		tx.Rollback() // 回滚事务
		return http_error_code.BadRequest("被合并账本不存在本人名下",
			http_error_code.WithInternal(err))
	}
	// 更新账本数值
	var transactions []model.Transaction
	if err := tx.Model(&model.Transaction{}).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("account_book_id = ?", mergeAccountBookReq.MergedAccountBookID).Find(&transactions).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	for _, userID := range mergedAccountBook.UserIDs {
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
