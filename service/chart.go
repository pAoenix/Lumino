package service

import (
	"Lumino/common/http_error_code"
	"Lumino/model"
	"Lumino/store"
	"fmt"
	"time"
)

// ChartService -
type ChartService struct {
	TransactionStore *store.TransactionStore
}

// NewChartService -
func NewChartService(transactionStore *store.TransactionStore) *ChartService {
	return &ChartService{
		TransactionStore: transactionStore,
	}
}

// GetNormalChart -
func (s *ChartService) GetNormalChart(chartReq *model.GetTransactionReq) (resp model.ChartResp, err error) {
	// 如果不传时间范围，默认最近7天
	if chartReq.EndTime == nil {
		*chartReq.EndTime = time.Now().AddDate(0, 0, 1)
	}
	if chartReq.BeginTime == nil {
		*chartReq.BeginTime = time.Now().Truncate(24*time.Hour).AddDate(0, 0, -6)
	}
	if chartReq.EndTime.Before(*chartReq.BeginTime) {
		return resp, http_error_code.BadRequest("时间范围异常, end需要>begin")
	}
	transactions, err := s.TransactionStore.Get(chartReq)
	if err != nil {
		return
	}
	for i, transaction := range transactions {
		fmt.Println(i, transaction)
	}
	return
}
