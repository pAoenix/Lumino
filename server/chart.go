package server

import (
	"Lumino/model"
	"Lumino/router/middleware"
	"Lumino/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ChartServer -
type ChartServer struct {
	ChartService *service.ChartService
}

// NewChartServer -
func NewChartServer(chartService *service.ChartService) *ChartServer {
	return &ChartServer{
		ChartService: chartService,
	}
}

// GetNormalChart -
// @Summary	统计看板
// @Tags 图表
// @Param        transaction  query      model.GetTransactionReq  true  "交易信息"
// @Success	200 {object}  model.ChartResp "图表信息"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/chart [get]
func (s *ChartServer) GetNormalChart(c *gin.Context) {
	req := model.GetTransactionReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if resp, err := s.ChartService.GetNormalChart(&req); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}
