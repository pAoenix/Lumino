package server

import (
	"Lumino/common/http_error_code"
	"Lumino/model"
	"Lumino/router/middleware"
	"Lumino/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AccountBookServer -
type AccountBookServer struct {
	AccountBookService *service.AccountBookService
}

// NewAccountBookServer -
func NewAccountBookServer(accountBookService *service.AccountBookService) *AccountBookServer {
	return &AccountBookServer{
		AccountBookService: accountBookService,
	}
}

// Register -
// @Summary	注册账本
// @Tags 账本
// @Param        account_book  query      model.RegisterAccountBookReq  true  "账本信息"
// @Success	200 {object}  model.AccountBook             "注册结果"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/account-book [post]
func (s *AccountBookServer) Register(c *gin.Context) {
	req := model.RegisterAccountBookReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if resp, err := s.AccountBookService.Register(&req); err != nil {
		c.Error(err)
		return
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Merge -
// @Summary	合并账本
// @Tags 账本
// @Param        account_book  query      model.MergeAccountBookReq  true  "账本id信息"
// @Success	200 {object}  model.AccountBookResp "账本结果"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/account-book/merge [post]
func (s *AccountBookServer) Merge(c *gin.Context) {
	req := model.MergeAccountBookReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if req.MergedAccountBookID == req.MergeAccountBookID {
		c.Error(http_error_code.BadRequest("合并账本不能相同"))
		return
	}
	if resp, err := s.AccountBookService.Merge(&req); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Get -
// @Summary	获取账本
// @Tags 账本
// @Param        account_book  query      model.GetAccountBookReq  true  "账本id信息"
// @Success	200 {object}  model.AccountBookResp "账本结果"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/account-book [get]
func (s *AccountBookServer) Get(c *gin.Context) {
	req := model.GetAccountBookReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if resp, err := s.AccountBookService.Get(&req); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Modify -
// @Summary	修改账本
// @Tags 账本
// @Param        account_book  query      model.ModifyAccountBookReq  true  "账本信息"
// @Success	200 {object}  model.AccountBook             "账本结果"
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/account-book [put]
func (s *AccountBookServer) Modify(c *gin.Context) {
	req := model.ModifyAccountBookReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if resp, err := s.AccountBookService.Modify(&req); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Delete -
// @Summary	删除账本
// @Tags 账本
// @Param        account_book  query      model.DeleteAccountBookReq  true  "账本信息"
// @Success	204
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/account-book [delete]
func (s *AccountBookServer) Delete(c *gin.Context) {
	req := model.DeleteAccountBookReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if err := s.AccountBookService.Delete(&req); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

// AA -
// @Summary	AA分账
// @Tags 账本
// @Param        account_book  query      model.AAAccountBookReq  true  "账本信息"
// @Success	200 {object}  []model.AAResult
// @Failure	400 {object}  http_error_code.AppError      "请求体异常"
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/account-book [delete]
func (s *AccountBookServer) AA(c *gin.Context) {
	req := model.AAAccountBookReq{}
	if err := middleware.Bind(c, &req); err != nil {
		c.Error(err)
		return
	}
	if resp, err := s.AccountBookService.AA(&req); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}
