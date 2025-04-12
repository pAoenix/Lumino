package server

import (
	"Lumino/common"
	"Lumino/common/http_error_code"
	"Lumino/model"
	"Lumino/router/middleware"
	"Lumino/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserServer -
type UserServer struct {
	UserService *service.UserService
}

// NewUserServer -
func NewUserServer(userService *service.UserService, client *common.OssClient) *UserServer {
	return &UserServer{
		UserService: userService,
	}
}

// Register -
// @Summary	注册用户信息
// @Tags 用户
// @Param        user  query      model.RegisterUserReq  true  "用户信息"
// @Param        icon_file formData file true "用户头像"
// @Success	200 {object}  model.User "注册结果"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/user [post]
func (s *UserServer) Register(c *gin.Context) {
	req := model.RegisterUserReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	// 1. 获取上传的文件
	fileHeader, err := c.FormFile("icon_file")
	if err != nil {
		c.Error(http_error_code.BadRequest("需要注册头像",
			http_error_code.WithInternal(err)))
		return
	}

	// 2. 验证图片类型
	ok, _, err := common.IsValidImage(fileHeader)
	if err != nil {
		c.Error(http_error_code.Internal("文件检查失败",
			http_error_code.WithInternal(err)))
		return
	}
	if !ok {
		c.Error(http_error_code.BadRequest("仅支持JPEG/PNG/GIF/BMP/WEBP图片",
			http_error_code.WithInternal(err)))
		return
	}
	// 注册入库
	resp, err := s.UserService.Register(&req, fileHeader)
	if err != nil {
		c.Error(err) // 交给中间件处理
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

// Modify -
// @Summary	修改用户信息
// @Tags 用户
// @Param        user    query      model.ModifyUserReq  true  "用户信息"
// @Param balance_detail query object false "余额详情" collectionFormat: multi
// @Success	200 {object}  model.User "用户修改后结果"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/user [put]
func (s *UserServer) Modify(c *gin.Context) {
	req := model.ModifyUserReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if resp, err := s.UserService.Modify(&req); err != nil {
		c.Error(err)
		return
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// ModifyProfilePhoto -
// @Summary	修改用户头像
// @Tags 用户
// @Param        user    query      model.ModifyProfilePhotoReq  true  "用户信息"
// @Param        icon_file formData file false "用户头像"
// @Success	204
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/user/profile-photo [put]
func (s *UserServer) ModifyProfilePhoto(c *gin.Context) {
	req := model.ModifyProfilePhotoReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	// 获取上传的文件
	fileHeader, err := c.FormFile("icon_file")
	if err != nil {
		c.Error(http_error_code.BadRequest("需要有头像文件",
			http_error_code.WithInternal(err)))
		return
	}
	// 2. 验证图片类型
	ok, _, err := common.IsValidImage(fileHeader)
	if err != nil {
		c.Error(http_error_code.Internal("文件检查失败",
			http_error_code.WithInternal(err)))
		return
	}
	if !ok {
		c.Error(http_error_code.BadRequest("仅支持JPEG/PNG/GIF/BMP/WEBP图片",
			http_error_code.WithInternal(err)))
		return
	}
	// 注册入库
	err = s.UserService.ModifyProfilePhoto(&req, fileHeader)
	if err != nil {
		c.Error(err) // 交给中间件处理
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

// Get -
// @Summary	获取用户信息
// @Tags 用户
// @Param        user  query      model.GetUserReq  true  "用户信息"
// @Success	200 {object}  model.User "用户结果"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/user [get]
func (s *UserServer) Get(c *gin.Context) {
	req := model.GetUserReq{}
	if err := middleware.BindQuery(c, &req); err != nil {
		c.Error(err)
		return
	}
	if user, err := s.UserService.Get(&req); err != nil {
		c.Error(err) // 交给中间件处理
		return
	} else {
		c.JSON(http.StatusOK, user)
	}
	return
}

// Delete -
// @Summary	删除用户信息
// @Tags 用户
// @Param        user  query      model.DeleteUserReq  true  "用户信息"
// @Success	204
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/user [delete]
func (s *UserServer) Delete(c *gin.Context) {
	req := model.DeleteUserReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if err := s.UserService.Delete(&req); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}
