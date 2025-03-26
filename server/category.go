package server

import (
	"Lumino/model"
	"Lumino/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CategoryServer -
type CategoryServer struct {
	CategoryService *service.CategoryService
}

// NewCategoryServer -
func NewCategoryServer(CategoryService *service.CategoryService) *CategoryServer {
	return &CategoryServer{
		CategoryService: CategoryService,
	}
}

// Register 注册图标
// @Summary	注册图标
// @Tags 图标
// @Param        category  body      model.Category  true  "图标信息"
// @Success	204
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/category [post]
func (s *CategoryServer) Register(c *gin.Context) {
	req := model.Category{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.CategoryService.Register(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

// Get -
// @Summary	获取图标
// @Tags 图标
// @Param        category  query      model.CategoryReq  true  "图标请求"
// @Success	200  {object} []model.Category "图标结果"
// @Failure	400  {string}  string      "请求体异常"
// @Failure	500  {string}  string      "服务端异常"
// @Router		/api/v1/category [get]
func (s *CategoryServer) Get(c *gin.Context) {
	req := model.CategoryReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if resp, err := s.CategoryService.Get(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Modify -
// @Summary	修改图标
// @Tags 图标
// @Param        category  body      model.Category  true  "图标信息"
// @Success	204
// @Failure	400  {string}  string      "请求体异常"
// @Failure	500  {string}  string      "服务端异常"
// @Router		/api/v1/category [put]
func (s *CategoryServer) Modify(c *gin.Context) {
	req := model.Category{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.CategoryService.Modify(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}
