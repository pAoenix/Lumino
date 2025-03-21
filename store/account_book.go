package store

import "Lumino/model"

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
	if s.db.Model(&model.AccountBook{}).Order(sort).Where("? = any(user_id)", accountBookReq.UserId).Find(&resp).Error != nil {
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
	return s.db.Model(&model.Transaction{}).
		Where("account_book_id = ?", mergeAccountBookReq.MergedAccountBookID).
		Update("account_book_id", mergeAccountBookReq.MergeAccountBookID).Error
}
