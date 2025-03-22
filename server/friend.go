package server

import (
	"Lumino/model"
	"Lumino/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// FriendServer -
type FriendServer struct {
	FriendService *service.FriendService
}

// NewFriendServer -
func NewFriendServer(friendService *service.FriendService) *FriendServer {
	return &FriendServer{
		FriendService: friendService,
	}
}

// Invite -
func (s *FriendServer) Invite(c *gin.Context) {
	req := model.Friend{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.FriendService.Invite(&req); err != nil {
		if err.Error() == "你已存在该好友" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

// Delete -
func (s *FriendServer) Delete(c *gin.Context) {
	req := model.Friend{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.FriendService.Delete(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}
