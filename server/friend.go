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
// @Summary	邀请朋友
// @Tags 朋友
// @Param        friend  query      model.Friend  true  "邀请信息"
// @Success	204
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/friend/invite [post]
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
// @Summary	删除朋友
// @Tags 朋友
// @Param        friend  query      model.Friend  true  "删除信息"
// @Success	204
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/friend [delete]
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
