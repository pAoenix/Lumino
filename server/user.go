package server

import (
	"Lumino/common"
	"Lumino/common/http_error_code"
	"Lumino/model"
	"Lumino/router/middleware"
	"Lumino/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
)

// UserServer -
type UserServer struct {
	UserService *service.UserService
	OssClient   *common.OssClient
}

// NewUserServer -
func NewUserServer(userService *service.UserService, client *common.OssClient) *UserServer {
	return &UserServer{
		UserService: userService,
		OssClient:   client,
	}
}

// Register -
// @Summary	注册用户信息
// @Tags 用户
// @Param        user  query      model.RegisterUserReq  true  "用户信息"
// @Param balance_detail query object false "余额详情" collectionFormat: multi
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
	// 获取上传的文件
	fileHeader, err := c.FormFile("icon_file")
	if err != nil {
		c.Error(http_error_code.BadRequest("需要注册头像",
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
// @Param        icon_file formData file false "用户头像"
// @Success	200 {object}  model.User "用户修改后结果"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/user [put]
func (s *UserServer) Modify(c *gin.Context) {
	req := model.ModifyUserReq{}
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
		req.IconUrl = viper.GetString("oss.profilePhotoDir") + strconv.Itoa(int(req.ID)) + ".jpg"
		if err := s.OssClient.UploadFile(req.IconUrl, file); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	if resp, err := s.UserService.Modify(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, resp)
	}
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
		if ossUrl, err := s.OssClient.DownloadFile(user.IconUrl); err != nil {
			c.Error(err) // 交给中间件处理
			return
		} else {
			user.IconUrl = ossUrl
		}
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
