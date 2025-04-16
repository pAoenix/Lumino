package server

import (
	"Lumino/model"
	"Lumino/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AccountServer -
type AccountServer struct {
	AccountService *service.AccountService
}

// NewAccountServer -
func NewAccountServer(AccountService *service.AccountService) *AccountServer {
	return &AccountServer{
		AccountService: AccountService,
	}
}

// Register -
// @Summary	注册账户
// @Tags 账户
// @Param        account  query      model.Account  true  "账户信息"
// @Success	204
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/account [post]
func (s *AccountServer) Register(c *gin.Context) {
	req := model.Account{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.AccountService.Register(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

// Modify -
// @Summary	修改账户
// @Tags 账户
// @Param        account  query      model.Account  true  "账户信息"
// @Success	204
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/account [put]
func (s *AccountServer) Modify(c *gin.Context) {
	req := model.Account{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.AccountService.Modify(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

// Get -
// @Summary	获取账户
// @Tags 账户
// @Param        account  query      model.Account  true  "账户信息"
// @Success	200 {object}  []model.Account "账户结果"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/account [get]
func (s *AccountServer) Get(c *gin.Context) {
	req := model.Account{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if resp, err := s.AccountService.Get(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Delete -
// @Summary	删除账户
// @Tags 账户
// @Param        account  query      model.Account  true  "账户信息"
// @Success	204
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/account [delete]
func (s *AccountServer) Delete(c *gin.Context) {
	req := model.Account{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.AccountService.Delete(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}
