package service

import (
	"Lumino/common"
	"Lumino/model"
	"Lumino/store"
	"time"
)

// TransactionService -
type TransactionService struct {
	TransactionStore   *store.TransactionStore
	userDownloader     UserIconDownloader
	categoryDownloader CategoryDownloader
}

// NewTransactionService -
func NewTransactionService(transactionStore *store.TransactionStore, userDownloader UserIconDownloader, categoryDownloader CategoryDownloader) *TransactionService {
	return &TransactionService{
		TransactionStore:   transactionStore,
		userDownloader:     userDownloader,
		categoryDownloader: categoryDownloader,
	}
}

// Register -
func (s *TransactionService) Register(transactionReq *model.RegisterTransactionReq) (resp model.Transaction, err error) {
	if transactionReq.Date.IsZero() {
		transactionReq.Date = time.Now()
	}
	return s.TransactionStore.Register(transactionReq)
}

// Get -
func (s *TransactionService) Get(transactionReq *model.GetTransactionReq) (resp model.TransactionResp, err error) {
	transactions, err := s.TransactionStore.Get(transactionReq)
	if err != nil {
		return
	}
	resp.Transactions = transactions
	// 获取全量用户信息
	var userIDs []uint
	var categoryIDs []uint
	for _, transaction := range transactions {
		if !common.ContainsUint(userIDs, transaction.PayUserID) {
			userIDs = append(userIDs, transaction.PayUserID)
		}
		if !common.ContainsUint(userIDs, transaction.CreatorID) {
			userIDs = append(userIDs, transaction.CreatorID)
		}
		for _, userID := range transaction.RelatedUserIDs {
			if !common.ContainsUint(userIDs, uint(userID)) {
				userIDs = append(userIDs, uint(userID))
			}
		}
		if !common.ContainsUint(categoryIDs, transaction.CategoryID) {
			categoryIDs = append(categoryIDs, transaction.CategoryID)
		}
	}
	if users, err := s.userDownloader.DownloadUserIcons(userIDs); err != nil {
		return resp, err
	} else {
		resp.Users = users
	}
	if categorys, err := s.categoryDownloader.DownloadCategoryIcon(categoryIDs, nil); err != nil {
		return resp, err
	} else {
		resp.Categorys = categorys
	}
	return
}

// Modify -
func (s *TransactionService) Modify(transactionReq *model.ModifyTransactionReq) (transaction model.Transaction, err error) {
	return s.TransactionStore.Modify(transactionReq)
}

// Delete -
func (s *TransactionService) Delete(transactionReq *model.DeleteTransactionReq) error {
	return s.TransactionStore.Delete(transactionReq)
}
