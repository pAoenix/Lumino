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

// Register 注册新的图标类型
// @Summary      注册新的图标类型
// tags category
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /api/v1/register [post]
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
