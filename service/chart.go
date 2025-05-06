package service

import (
	"Lumino/common"
	"Lumino/common/http_error_code"
	"Lumino/model"
	"Lumino/store"
	"sort"
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
		endTime := time.Now()
		chartReq.EndTime = &endTime
	}
	if chartReq.BeginTime == nil {
		now := time.Now()
		beginTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -6)
		chartReq.BeginTime = &beginTime
	}
	if chartReq.EndTime.Before(*chartReq.BeginTime) {
		return resp, http_error_code.BadRequest("时间范围异常, end需要>begin")
	}
	transactions, err := s.TransactionStore.Get(chartReq)
	if err != nil {
		return
	}
	period := model.DayPeriod
	subDays := chartReq.EndTime.Sub(*chartReq.BeginTime).Hours() / 24
	if subDays > 31 {
		period = model.MonthPeriod
	}
	categoryMap := map[uint][]model.Transaction{}
	dateMap := map[string][]model.Transaction{}
	for _, transaction := range transactions {
		resp.TotalAmount += transaction.Amount
		// 按照时间管理
		currentDate := transaction.Date.Format("2006-01-02")
		if period == model.MonthPeriod {
			currentDate = transaction.Date.Format("2006-01")
		}
		if _, ok := dateMap[currentDate]; !ok {
			dateMap[currentDate] = []model.Transaction{transaction}
		} else {
			dateMap[currentDate] = append(dateMap[currentDate], transaction)
		}
		// 按照类型管理
		if _, ok := categoryMap[transaction.CategoryID]; !ok {
			categoryMap[transaction.CategoryID] = []model.Transaction{transaction}
		} else {
			categoryMap[transaction.CategoryID] = append(categoryMap[transaction.CategoryID], transaction)
		}
	}
	// 计算好每个类型的花销数据
	for key, value := range categoryMap {
		amount := 0.
		sort.Slice(value, func(i, j int) bool {
			return value[i].Amount > value[j].Amount
		})
		for _, v := range value {
			amount += v.Amount
		}
		resp.CategoryChart = append(resp.CategoryChart, model.CategoryChart{
			CategoryID:   key,
			Amount:       amount,
			Percent:      amount / resp.TotalAmount,
			Transactions: value,
		})
	}
	sort.Slice(resp.CategoryChart, func(i, j int) bool {
		return resp.CategoryChart[i].Amount > resp.CategoryChart[j].Amount
	})
	// 计算好每个时间点的花销数据
	var dateList []string
	for key, value := range dateMap {
		amount := 0.
		dateList = append(dateList, key)
		sort.Slice(value, func(i, j int) bool {
			return value[i].Amount > value[j].Amount
		})
		for _, v := range value {
			amount += v.Amount
		}
		resp.DateChart = append(resp.DateChart, model.DateChart{
			DateStr:      key,
			Amount:       amount,
			Transactions: value,
		})
	}
	if period == model.DayPeriod {
		if ok, addDate := common.CheckDailyCoverage(*chartReq.BeginTime, *chartReq.EndTime, dateList); !ok {
			for _, ad := range addDate {
				resp.DateChart = append(resp.DateChart, model.DateChart{
					DateStr:      ad,
					Amount:       0.,
					Transactions: []model.Transaction{},
				})
			}
		}
	} else if period == model.MonthPeriod {
		if ok, addDate := common.CheckMonthlyCoverage(*chartReq.BeginTime, *chartReq.EndTime, dateList); !ok {
			for _, ad := range addDate {
				resp.DateChart = append(resp.DateChart, model.DateChart{
					DateStr:      ad,
					Amount:       0.,
					Transactions: []model.Transaction{},
				})
			}
		}
	}
	sort.Slice(resp.DateChart, func(i, j int) bool {
		return resp.DateChart[i].DateStr < resp.DateChart[j].DateStr
	})
	if len(resp.DateChart) == 0 {
		resp.AverageAmount = 0
	} else {
		resp.AverageAmount = resp.TotalAmount / float64(len(resp.DateChart))
	}
	return
}
