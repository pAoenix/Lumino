package server

import (
	"Lumino/common"
	"Lumino/model"
	"Lumino/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"sync"
)

// CategoryServer -
type CategoryServer struct {
	CategoryService *service.CategoryService
	OssClient       *common.OssClient
}

// NewCategoryServer -
func NewCategoryServer(CategoryService *service.CategoryService, client *common.OssClient) *CategoryServer {
	return &CategoryServer{
		CategoryService: CategoryService,
		OssClient:       client,
	}
}

// Register 注册图标
// @Summary	注册图标
// @Tags 图标
// @Param        category  query      model.RegisterCategoryReq  true  "图标信息"
// @Param        icon_file formData file true "分类图标文件"
// @Success	204
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/category [post]
func (s *CategoryServer) Register(c *gin.Context) {
	req := model.RegisterCategoryReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// 获取上传的文件
	fileHeader, err := c.FormFile("icon_file")
	if err != nil {
		c.JSON(400, gin.H{"error": "需要注册分类图标文件"})
		return
	}
	// 2. 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打开文件失败: " + err.Error()})
		return
	}
	// 3. 注册入库
	resp, err := s.CategoryService.Register(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	// 4. 上传文件
	defer file.Close()
	iconUrl := viper.GetString("oss.categoryDir") + strconv.Itoa(int(resp.ID)) + ".jpg"
	if err := s.OssClient.UploadFile(iconUrl, file); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	// 5. 更新icon地址
	modifyReq := model.ModifyCategoryReq{ID: resp.ID, IconUrl: iconUrl}
	if resp, err := s.CategoryService.Modify(&modifyReq); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if resp, err := s.CategoryService.Get(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		const maxConcurrency = 20 // 最大并发数
		sem := make(chan struct{}, maxConcurrency)
		var wg sync.WaitGroup
		for idx, _ := range resp {
			wg.Add(1) // 计数器加1
			go func(i int) {
				sem <- struct{}{} // 获取信号量
				defer func() {
					<-sem // 释放信号量
					wg.Done()
				}()
				if ossUrl, err := s.OssClient.DownloadFile(resp[i].IconUrl); err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				} else {
					resp[i].IconUrl = ossUrl
				}
			}(idx)
		}
		wg.Wait() // 等待所有goroutine完成
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Modify -
// @Summary	修改图标
// @Tags 图标
// @Param        category  query      model.ModifyCategoryReq  true  "图标信息"
// @Param        icon_file formData file false "分类图标文件"
// @Success	204
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/category [put]
func (s *CategoryServer) Modify(c *gin.Context) {
	req := model.ModifyCategoryReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// 获取上传的文件
	fileHeader, err := c.FormFile("icon_file")
	// 有文件才进行头像修改
	if err == nil {
		// 2. 打开文件
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "打开文件失败: " + err.Error()})
			return
		}
		// 3. 上传文件
		defer file.Close()
		req.IconUrl = viper.GetString("oss.categoryDir") + strconv.Itoa(int(req.ID)) + ".jpg"
		if err := s.OssClient.UploadFile(req.IconUrl, file); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	if resp, err := s.CategoryService.Modify(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	} else {
		c.JSON(http.StatusNoContent, resp)
	}
	return
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
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.CategoryService.Delete(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}
