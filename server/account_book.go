package server

import (
	"Lumino/common"
	"Lumino/model"
	"Lumino/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AccountBookServer -
type AccountBookServer struct {
	AccountBookService *service.AccountBookService
}

// NewAccountBookServer -
func NewAccountBookServer(accountBookService *service.AccountBookService) *AccountBookServer {
	return &AccountBookServer{
		AccountBookService: accountBookService,
	}
}

// Register -
func (s *AccountBookServer) Register(c *gin.Context) {
	req := model.AccountBook{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.AccountBookService.Register(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

// Merge -
func (s *AccountBookServer) Merge(c *gin.Context) {
	req := model.MergeAccountBookReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.AccountBookService.Merge(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

// Get -
func (s *AccountBookServer) Get(c *gin.Context) {
	req := model.AccountBookReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if resp, err := s.AccountBookService.Get(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// GetByID -
func (s *AccountBookServer) GetByID(c *gin.Context) {
	accountBookID := c.Param("id")
	abID, _ := common.String2Uint(accountBookID)
	req := model.AccountBookReq{ID: abID}
	if resp, err := s.AccountBookService.Get(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Modify -
func (s *AccountBookServer) Modify(c *gin.Context) {
	req := model.AccountBook{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.AccountBookService.Modify(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

// Delete -
func (s *AccountBookServer) Delete(c *gin.Context) {
	req := model.AccountBook{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.AccountBookService.Delete(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}
