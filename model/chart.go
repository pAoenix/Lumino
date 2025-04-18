package model

import "time"

// ChartReq 获取图表请求
type ChartReq struct {
	UserID     uint       `json:"user_id" form:"user_id" binding:"required"` // 用户id
	BeginTime  *time.Time `json:"begin_time" form:"begin_time"`              // 起始时间
	EndTime    *time.Time `json:"end_time" form:"end_time"`                  // 结束时间
	Type       int        `json:"type" form:"type"`                          // 类型:收入/支出
	CategoryID uint       `json:"category_id" form:"category_id"`            // 类别
	Period     string     `json:"period" form:"Period"`                      // 周期:周，月，年
}

// CategoryChart -
type CategoryChart struct {
	CategoryID  uint          `json:"category_id" form:"category_id"` // 关联消费场景分类ID
	Amount      float64       `json:"amount" form:"amount"`           // 交易数额
	Percent     float64       `json:"percent" form:"percent"`         // 占比
	Transaction []Transaction `json:"transaction" form:"transaction"` // 交易记录,全部的
}

// DateChart -
type DateChart struct {
	DateStr     string        `json:"date_str" form:"date_str"`       // 报表粒度，月 or 天
	Amount      float64       `json:"amount" form:"amount"`           // 交易数额
	Transaction []Transaction `json:"transaction" form:"transaction"` // 交易记录，top3
}

// ChartResp -
type ChartResp struct {
	TotalAmount   float64         `json:"total_amount" form:"total_amount"`     // 总费用
	AverageAmount float64         `json:"average_amount" form:"average_amount"` // 平均费用
	DateChart     []DateChart     `json:"date_chart" form:"date_chart"`         // 时间详情
	CategoryChart []CategoryChart `json:"category_chart" form:"category_chart"` // 分类详情
}
