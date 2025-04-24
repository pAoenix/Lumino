package service

import (
	"Lumino/common"
	"Lumino/common/http_error_code"
	"Lumino/model"
	"Lumino/store"
	"sort"
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
	// 如果没有时间范围，那么默认为本月
	if transactionReq.BeginTime == nil {
		firstDay := common.GetFirstDayOfMonth(time.Now())
		transactionReq.BeginTime = &firstDay
	}
	if transactionReq.EndTime == nil {
		LastDay := common.GetLastDayOfMonth(time.Now()).AddDate(0, 0, 1)
		transactionReq.EndTime = &LastDay
	}
	if transactionReq.EndTime.Before(*transactionReq.BeginTime) {
		return resp, http_error_code.BadRequest("时间范围异常, end需要>begin")
	}
	transactions, err := s.TransactionStore.Get(transactionReq)
	if err != nil {
		return
	}
	resp.Transactions = groupByDay(transactions)
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

// 按天分组记账数据
func groupByDay(items []model.Transaction) (dailyTransaction []model.DailyTransaction) {
	// 按天分组
	dailyMap := make(map[string]model.DailyTransaction)
	var keys []string
	for _, item := range items {
		dateDayStr := item.Date.Format("2006-01-02")
		if daily, exists := dailyMap[dateDayStr]; exists {
			daily.Items = append(daily.Items, item)
			if item.Type == model.IncomeType {
				daily.Income += item.Amount
			} else {
				daily.Spending += item.Amount
			}
			dailyMap[dateDayStr] = daily
		} else {
			spending := 0.
			income := 0.
			if item.Type == model.IncomeType {
				income = item.Amount
			} else {
				spending = item.Amount
			}
			keys = append(keys, dateDayStr)
			dailyMap[dateDayStr] = model.DailyTransaction{
				Date:     dateDayStr,
				Items:    []model.Transaction{item},
				Spending: spending,
				Income:   income,
			}
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	// 将map转换为slice
	for _, key := range keys {
		dailyTransaction = append(dailyTransaction, dailyMap[key])
	}
	return
}
