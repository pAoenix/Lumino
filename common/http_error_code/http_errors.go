package http_error_code

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

// ErrorType 定义错误分类
type ErrorType string

const (
	ErrorTypeInvalidInput ErrorType = "invalid_input" // 400
	ErrorTypeUnauthorized ErrorType = "unauthorized"  // 401
	ErrorTypeForbidden    ErrorType = "forbidden"     // 403
	ErrorTypeNotFound     ErrorType = "not_found"     // 404
	ErrorTypeConflict     ErrorType = "conflict"      // 409
	ErrorTypeInternal     ErrorType = "internal"      // 500
	ErrorTypeUnavailable  ErrorType = "unavailable"   // 503
)

// AppError 应用错误结构体
type AppError struct {
	Type     ErrorType `json:"type"`             // 错误类型
	Code     int       `json:"code"`             // HTTP状态码
	Message  string    `json:"message"`          // 用户友好消息
	Detail   string    `json:"detail,omitempty"` // 调试细节
	Internal error     `json:"-"`                // 内部错误(不暴露给客户端)
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	return e.Message
}

// Unwrap 支持errors.Unwrap
func (e *AppError) Unwrap() error {
	return e.Internal
}

// New 创建新错误
func New(typ ErrorType, msg string, detail string, err error) *AppError {
	code := typeToCode(typ)
	return &AppError{
		Type:     typ,
		Code:     code,
		Message:  msg,
		Detail:   detail,
		Internal: err,
	}
}

// 预定义错误构造器

// BadRequest -
func BadRequest(msg, detail string, err error) *AppError {
	return New(ErrorTypeInvalidInput, msg, detail, err)
}

// Unauthorized -
func Unauthorized(msg, detail string, err error) *AppError {
	return New(ErrorTypeUnauthorized, msg, detail, err)
}

// NotFound -
func NotFound(msg, detail string, err error) *AppError {
	return New(ErrorTypeNotFound, msg, detail, err)
}

// InternalServer -
func InternalServer(msg, detail string, err error) *AppError {
	return New(ErrorTypeInternal, msg, detail, err)
}

// FromDB 从GORM错误转换
func FromDB(err error) *AppError {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return BadRequest("请求内容不存在", "", err)
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return New(ErrorTypeConflict, "数据冲突", "", err)
	case errors.Is(err, gorm.ErrInvalidTransaction):
		return InternalServer("数据库事务错误", "", err)
	default:
		return InternalServer("数据库操作失败", "", err)
	}
}

// typeToCode 错误类型到HTTP状态码
func typeToCode(typ ErrorType) int {
	switch typ {
	case ErrorTypeInvalidInput:
		return http.StatusBadRequest
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeInternal:
		return http.StatusInternalServerError
	case ErrorTypeUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
