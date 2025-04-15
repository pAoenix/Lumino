package store

import (
	"Lumino/common/http_error_code"
	"Lumino/model"
	"github.com/jinzhu/copier"
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
func (s *TransactionStore) Register(transactionReq *model.RegisterTransactionReq) (transaction model.Transaction, err error) {
	if err = ParamsJudge(s.db, transactionReq.AccountBookID, transactionReq.RelatedUserIDs,
		transactionReq.CreatorID, transactionReq.CategoryID, nil); err != nil {
		return transaction, err
	}
	if err = copier.Copy(&transaction, &transactionReq); err != nil {
		return transaction, http_error_code.Internal("服务内异常",
			http_error_code.WithInternal(err))
	}
	if err = s.db.Model(&model.Transaction{}).Create(&transaction).Error; err != nil {
		return
	}
	return
}

// Get -
func (s *TransactionStore) Get(transactionReq *model.GetTransactionReq) (resp []model.Transaction, err error) {
	if err = ParamsJudge(s.db, &transactionReq.AccountBookID, nil,
		nil, nil, transactionReq.ID); err != nil {
		return resp, err
	}
	sql := s.db.Model(&model.Transaction{})
	if transactionReq.BeginTime != nil {
		sql.Where("date >= ?", &transactionReq.BeginTime)
	}
	if transactionReq.EndTime != nil {
		sql.Where("date <= ?", &transactionReq.EndTime)
	}
	if transactionReq.ID != nil {
		sql.Where("id = ？", &transactionReq.ID)
	}
	err = sql.Where("account_book_id = ?", transactionReq.AccountBookID).Find(&resp).Error
	return
}

// Modify -
func (s *TransactionStore) Modify(modifyTransaction *model.ModifyTransactionReq) (transaction model.Transaction, err error) {
	if err = ParamsJudge(s.db, modifyTransaction.AccountBookID, modifyTransaction.RelatedUserIDs,
		modifyTransaction.PayUserID, modifyTransaction.CategoryID, &modifyTransaction.ID); err != nil {
		return
	}
	// 交易信息更新
	if err = copier.Copy(&transaction, &modifyTransaction); err != nil {
		return transaction, http_error_code.Internal("服务内异常",
			http_error_code.WithInternal(err))
	}
	err = s.db.Model(&model.Transaction{}).Where("id = ?", modifyTransaction.ID).Updates(&transaction).Find(&transaction).Error
	return
}

// Delete -
func (s *TransactionStore) Delete(transactionReq *model.DeleteTransactionReq) error {
	if err := ParamsJudge(s.db, nil, nil, nil, nil, &transactionReq.ID); err != nil {
		return err
	}
	return s.db.Model(&model.Transaction{}).Where("id = ?", transactionReq.ID).Delete(&model.Transaction{}).Error
}
