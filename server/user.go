package server

import (
	"Lumino/model"
	"Lumino/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserServer -
type UserServer struct {
	UserService *service.UserService
}

// NewUserServer -
func NewUserServer(UserService *service.UserService) *UserServer {
	return &UserServer{
		UserService: UserService,
	}
}

// Register -
func (s *UserServer) Register(c *gin.Context) {
	req := model.User{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.UserService.Register(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}
