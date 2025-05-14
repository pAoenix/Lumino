package server

import (
	"Lumino/model"
	"Lumino/router/middleware"
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
// @Param        account  query      model.RegisterAccountReq  true  "账户信息"
// @Success	200 {object}  model.Account "注册结果"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/account [post]
func (s *AccountServer) Register(c *gin.Context) {
	req := model.RegisterAccountReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if resp, err := s.AccountService.Register(&req); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Modify -
// @Summary	修改账户
// @Tags 账户
// @Param        account  query      model.ModifyAccountReq  true  "账户信息"
// @Success	200 {object}  model.Account "修改结果"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/account [put]
func (s *AccountServer) Modify(c *gin.Context) {
	req := model.ModifyAccountReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if resp, err := s.AccountService.Modify(&req); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Get -
// @Summary	获取账户
// @Tags 账户
// @Param       account  query      model.GetAccountReq  true  "账户信息"
// @Success	200 {object}  []model.Account "账户结果"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/account [get]
func (s *AccountServer) Get(c *gin.Context) {
	req := model.GetAccountReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if resp, err := s.AccountService.Get(&req); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Delete -
// @Summary	删除账户
// @Tags 账户
// @Param        account  query      model.DeleteAccountReq  true  "账户信息"
// @Success	204
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/account [delete]
func (s *AccountServer) Delete(c *gin.Context) {
	req := model.DeleteAccountReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if err := s.AccountService.Delete(&req); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}
