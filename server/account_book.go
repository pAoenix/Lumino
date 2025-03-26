package server

import (
	"Lumino/common"
	"Lumino/model"
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
// @Param        account_book  body      model.AccountBook  true  "账本信息"
// @Success	204
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/account-book [post]
func (s *AccountBookServer) Register(c *gin.Context) {
	req := model.AccountBook{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.AccountBookService.Register(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

// Merge -
// @Summary	合并账本
// @Tags 账本
// @Param        account_book  body      model.MergeAccountBookReq  true  "账本id信息"
// @Success	204
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/account-book/merge [post]
func (s *AccountBookServer) Merge(c *gin.Context) {
	req := model.MergeAccountBookReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.AccountBookService.Merge(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

// Get -
// @Summary	获取账本
// @Tags 账本
// @Param        account_book  query      model.AccountBookReq  true  "账本id信息"
// @Success	200 {object}  model.AccountBookResp "账本结果"
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/account-book [get]
func (s *AccountBookServer) Get(c *gin.Context) {
	req := model.AccountBookReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if resp, err := s.AccountBookService.Get(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// GetByID -
// @Summary	获取账本
// @Tags 账本
// @Param        id  path      int  true  "账本id" format(uint)
// @Success	200 {object}  model.AccountBookResp "账本结果"
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/account-book/:id [get]
func (s *AccountBookServer) GetByID(c *gin.Context) {
	accountBookID := c.Param("id")
	abID, _ := common.String2Uint(accountBookID)
	req := model.AccountBookReq{ID: abID}
	if resp, err := s.AccountBookService.Get(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, resp)
	}
	return
}

// Modify -
// @Summary	修改账本
// @Tags 账本
// @Param        account_book  body      model.AccountBook  true  "账本信息"
// @Success	204
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/account-book [put]
func (s *AccountBookServer) Modify(c *gin.Context) {
	req := model.AccountBook{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.AccountBookService.Modify(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

// Delete -
// @Summary	删除账本
// @Tags 账本
// @Param        account_book  body      model.AccountBook  true  "账本信息"
// @Success	204
// @Failure	400 {string}  string      "请求体异常"
// @Failure	500 {string}  string      "服务端异常"
// @Router		/api/v1/account-book [delete]
func (s *AccountBookServer) Delete(c *gin.Context) {
	req := model.AccountBook{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := s.AccountBookService.Delete(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}
