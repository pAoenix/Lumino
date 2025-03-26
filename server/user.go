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
// @Summary	注册用户信息
// @Tags 用户
// @Param        user  body      model.User  true  "用户信息"
// @Success	200 {object}  model.User "注册结果"
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/user [post]
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
	c.JSON(http.StatusOK, req)
	return
}

// Modify -
// @Summary	修改用户信息
// @Tags 用户
// @Param        user  body      model.User  true  "用户信息"
// @Success	200 {object}  model.User "用户修改后结果"
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/user [put]
func (s *UserServer) Modify(c *gin.Context) {
	req := model.User{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.UserService.Modify(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
	return
}

// Get -
// @Summary	获取用户信息
// @Tags 用户
// @Param        user  query      model.User  true  "用户信息"
// @Success	200 {object}  []model.User "用户结果"
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/user [get]
func (s *UserServer) Get(c *gin.Context) {
	req := model.User{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if users, err := s.UserService.Get(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, users)
	}
	return
}
