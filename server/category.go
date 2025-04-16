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

// CategoryServer -
type CategoryServer struct {
	CategoryService *service.CategoryService
}

// NewCategoryServer -
func NewCategoryServer(CategoryService *service.CategoryService, client *common.OssClient) *CategoryServer {
	return &CategoryServer{
		CategoryService: CategoryService,
	}
}

// Register 注册图标
// @Summary	注册图标
// @Tags 图标
// @Param        category  query      model.RegisterCategoryReq  true  "图标信息"
// @Param        icon_file formData file true "分类图标文件"
// @Success	200 {object}  model.Category                "图标信息"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/category [post]
func (s *CategoryServer) Register(c *gin.Context) {
	req := model.RegisterCategoryReq{}
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
	// 3. 注册入库
	if resp, err := s.CategoryService.Register(&req, fileHeader); err != nil {
		c.Error(err)
		return
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Get -
// @Summary	获取图标
// @Tags 图标
// @Param        category  query      model.GetCategoryReq  true  "图标请求"
// @Success	200  {object} []model.Category "图标结果"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/category [get]
func (s *CategoryServer) Get(c *gin.Context) {
	req := model.GetCategoryReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if resp, err := s.CategoryService.Get(&req); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Modify -
// @Summary	修改图标信息
// @Tags 图标
// @Param        category  query      model.ModifyCategoryReq  true  "图标信息"
// @Success	200 {object}  model.Category "图标结果"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/category [put]
func (s *CategoryServer) Modify(c *gin.Context) {
	req := model.ModifyCategoryReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if resp, err := s.CategoryService.Modify(&req); err != nil {
		c.Error(err)
		return
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// ModifyIconImage -
// @Summary	修改图标文件
// @Tags 图标
// @Param        category  query      model.ModifyCategoryIconReq  true  "图标信息"
// @Param        icon_file formData file true "分类图标文件"
// @Success	204
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/category/icon-image [put]
func (s *CategoryServer) ModifyIconImage(c *gin.Context) {
	// 获取上传的文件
	req := model.ModifyCategoryIconReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	fileHeader, err := c.FormFile("icon_file")
	if err != nil {
		c.Error(http_error_code.BadRequest("需要有图标文件",
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
	err = s.CategoryService.ModifyProfilePhoto(&req, fileHeader)
	if err != nil {
		c.Error(err) // 交给中间件处理
		return
	}
	c.JSON(http.StatusNoContent, nil)

}

// Delete -
// @Summary	删除图标
// @Tags 图标
// @Param        category  query      model.DeleteCategoryReq  true  "图标信息"
// @Success	204
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/category [delete]
func (s *CategoryServer) Delete(c *gin.Context) {
	req := model.DeleteCategoryReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if err := s.CategoryService.Delete(&req); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}
