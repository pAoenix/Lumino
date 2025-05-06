package model

import "time"

const DayPeriod = 0
const WeekPeriod = 1
const MonthPeriod = 2

// ChartReq 获取图表请求
type ChartReq struct {
	ID            *uint      `json:"id" form:"id"`                                              // 交易id
	UserID        uint       `json:"user_id" form:"user_id" binding:"required"`                 // 用户id
	AccountBookID uint       `json:"account_book_id" form:"account_book_id" binding:"required"` // 对应的账本id
	Type          int        `json:"type" form:"type" binding:"required"`                       // 类型:收入/支出
	BeginTime     *time.Time `json:"begin_time" form:"begin_time"`                              // 起始时间
	EndTime       *time.Time `json:"end_time" form:"end_time"`                                  // 结束时间
	CategoryID    uint       `json:"category_id" form:"category_id"`                            // 类别
	// Period 周期:周，月，年;如果时间范围<1个月，按天输出，如果>1个月，按月输出
	Period string `json:"period" form:"period" swaggerignore:"true"`
}

// CategoryChart -
type CategoryChart struct {
	CategoryID   uint          `json:"category_id" form:"category_id"`   // 关联消费场景分类ID
	Amount       float64       `json:"amount" form:"amount"`             // 交易数额
	Percent      float64       `json:"percent" form:"percent"`           // 占比
	Transactions []Transaction `json:"transactions" form:"transactions"` // 交易记录,全部的
}

// DateChart -
type DateChart struct {
	DateStr      string        `json:"date_str" form:"date_str"`         // 报表粒度，月 or 天
	Amount       float64       `json:"amount" form:"amount"`             // 交易数额
	Transactions []Transaction `json:"transactions" form:"transactions"` // 交易记录，top3
}

// ChartResp -
type ChartResp struct {
	TotalAmount   float64         `json:"total_amount" form:"total_amount"`     // 总费用
	AverageAmount float64         `json:"average_amount" form:"average_amount"` // 平均费用
	DateChart     []DateChart     `json:"date_chart" form:"date_chart"`         // 时间详情
	CategoryChart []CategoryChart `json:"category_chart" form:"category_chart"` // 分类详情
}
