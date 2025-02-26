package server

import (
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
