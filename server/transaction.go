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
// @Param        transaction  body      model.Transaction  true  "交易信息"
// @Success	204
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/transaction [post]
func (s *TransactionServer) Register(c *gin.Context) {
	req := model.Transaction{}
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
// @Param        transaction  query      model.TransactionReq  true  "交易信息"
// @Success	200 {object}  []model.Transaction "交易记录"
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/transaction [get]
func (s *TransactionServer) Get(c *gin.Context) {
	req := model.TransactionReq{}
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
// @Param        transaction  body      model.Transaction  true  "交易信息"
// @Success	204
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/transaction [put]
func (s *TransactionServer) Modify(c *gin.Context) {
	req := model.Transaction{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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
// @Param        transaction  body      model.Transaction  true  "交易信息"
// @Success	204
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/transaction [delete]
func (s *TransactionServer) Delete(c *gin.Context) {
	req := model.Transaction{}
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
