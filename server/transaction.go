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
