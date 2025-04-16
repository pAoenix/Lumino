package server

import (
	"Lumino/common/http_error_code"
	"Lumino/model"
	"Lumino/router/middleware"
	"Lumino/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// TransactionServer -
type TransactionServer struct {
	TransactionService *service.TransactionService
}

// NewTransactionServer -
func NewTransactionServer(transactionService *service.TransactionService) *TransactionServer {
	return &TransactionServer{
		TransactionService: transactionService,
	}
}

// Register -
// @Summary	注册交易记录
// @Tags 交易记录
// @Param        transaction  query      model.RegisterTransactionReq  true  "交易信息"
// @Success	200 {object}  model.Transaction             "交易记录"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/transaction [post]
func (s *TransactionServer) Register(c *gin.Context) {
	req := model.RegisterTransactionReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if resp, err := s.TransactionService.Register(&req); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Get -
// @Summary	获取交易记录
// @Tags 交易记录
// @Param        transaction  query      model.GetTransactionReq  true  "交易信息"
// @Success	200 {object}  model.TransactionResp "交易记录"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/transaction [get]
func (s *TransactionServer) Get(c *gin.Context) {
	req := model.GetTransactionReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if resp, err := s.TransactionService.Get(&req); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Modify -
// @Summary	修改交易记录
// @Tags 交易记录
// @Param        transaction  query      model.ModifyTransactionReq  true  "交易信息"
// @Success	200 {object}  model.Transaction             "交易信息"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/transaction [put]
func (s *TransactionServer) Modify(c *gin.Context) {
	req := model.ModifyTransactionReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if req.Amount < 0.0 || req.Type != 0 && req.Amount <= 0.0 {
		c.Error(http_error_code.BadRequest("金额需要>0"))
		return
	}
	if req.Type != 0 && req.Type != model.IncomeType && req.Type != model.SpendingType {
		c.Error(http_error_code.BadRequest("交易类型必须是0或者1"))
		return
	}
	if resp, err := s.TransactionService.Modify(&req); err != nil {
		c.Error(err)
		return
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Delete -
// @Summary	删除交易记录
// @Tags 交易记录
// @Param        transaction  query      model.DeleteTransactionReq  true  "交易信息"
// @Success	204
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/transaction [delete]
func (s *TransactionServer) Delete(c *gin.Context) {
	req := model.DeleteTransactionReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if err := s.TransactionService.Delete(&req); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}
