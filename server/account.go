package server

import (
	"Lumino/model"
	"Lumino/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AccountServer -
type AccountServer struct {
	AccountService *service.AccountService
}

// NewAccountServer -
func NewAccountServer(AccountService *service.AccountService) *AccountServer {
	return &AccountServer{
		AccountService: AccountService,
	}
}

// Register -
func (s *AccountServer) Register(c *gin.Context) {
	req := model.Account{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.AccountService.Register(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

// Modify -
func (s *AccountServer) Modify(c *gin.Context) {
	req := model.Account{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.AccountService.Modify(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

// Get -
func (s *AccountServer) Get(c *gin.Context) {
	req := model.Account{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.AccountService.Get(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
	return
}

// Delete -
func (s *AccountServer) Delete(c *gin.Context) {
	req := model.Account{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.AccountService.Delete(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}
