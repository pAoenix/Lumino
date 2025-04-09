package server

import (
	"Lumino/model"
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
// @Success	204
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/transaction [post]
func (s *TransactionServer) Register(c *gin.Context) {
	req := model.RegisterTransactionReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.TransactionService.Register(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

// Get -
// @Summary	获取交易记录
// @Tags 交易记录
// @Param        transaction  query      model.GetTransactionReq  true  "交易信息"
// @Success	200 {object}  []model.Transaction "交易记录"
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/transaction [get]
func (s *TransactionServer) Get(c *gin.Context) {
	req := model.GetTransactionReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if resp, err := s.TransactionService.Get(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Modify -
// @Summary	修改交易记录
// @Tags 交易记录
// @Param        transaction  query      model.ModifyTransactionReq  true  "交易信息"
// @Success	204
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/transaction [put]
func (s *TransactionServer) Modify(c *gin.Context) {
	req := model.ModifyTransactionReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if req.Amount < 0.0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "金额需要>0"})
		return
	}
	if req.Type != 0 && req.Type != model.IncomeType && req.Type != model.SpendingType {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "交易类型必须是0或者1"})
		return
	}
	if err := s.TransactionService.Modify(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

// Delete -
// @Summary	删除交易记录
// @Tags 交易记录
// @Param        transaction  body      model.DeleteTransactionReq  true  "交易信息"
// @Success	204
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/transaction [delete]
func (s *TransactionServer) Delete(c *gin.Context) {
	req := model.DeleteTransactionReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.TransactionService.Delete(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}
