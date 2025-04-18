package service

import (
	"Lumino/store"
)

// ChartService -
type ChartService struct {
	chartStore *store.ChartStore
}

// NewChartService -
func NewChartService(chartStore *store.ChartStore) *ChartService {
	return &ChartService{
		chartStore: chartStore,
	}
}
