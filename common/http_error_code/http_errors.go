package http_error_code

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// ErrorType 定义错误分类
type ErrorType string

const (
	ErrorTypeInvalidInput   ErrorType = "invalid_input"   // 400
	ErrorTypeUnauthorized   ErrorType = "unauthorized"    // 401
	ErrorTypeForbidden      ErrorType = "forbidden"       // 403
	ErrorTypeNotFound       ErrorType = "not_found"       // 404
	ErrorTypeConflict       ErrorType = "conflict"        // 409
	ErrorTypeRateLimited    ErrorType = "rate_limited"    // 429
	ErrorTypeInternal       ErrorType = "internal"        // 500
	ErrorTypeNotImplemented ErrorType = "not_implemented" // 501
	ErrorTypeUnavailable    ErrorType = "unavailable"     // 503
	// ErrorTypeBindingFailed should-bind重写
	ErrorTypeBindingFailed ErrorType = "binding_failed" // 400
)

// AppError 增强版应用错误结构体
type AppError struct {
	Type        ErrorType `json:"type"`              // 错误类型
	Code        int       `json:"code"`              // HTTP状态码
	Message     string    `json:"message"`           // 用户友好消息
	Detail      string    `json:"detail,omitempty"`  // 调试细节
	Internal    error     `json:"-"`                 // 内部错误(不暴露)
	StackTrace  []string  `json:"-"`                 // 调用堆栈(仅开发环境)
	ServiceName string    `json:"service,omitempty"` // 服务名称
}

// Error 实现error接口
func (e *AppError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] %s", e.Type, e.Message))

	if e.Detail != "" {
		sb.WriteString(fmt.Sprintf(" (detail: %s)", e.Detail))
	}

	if e.Internal != nil {
		sb.WriteString(fmt.Sprintf(": %v", e.Internal))
	}

	return sb.String()
}

// Unwrap 支持errors.Unwrap
func (e *AppError) Unwrap() error {
	return e.Internal
}

// WithStack 添加调用堆栈
func (e *AppError) WithStack() *AppError {
	if e.StackTrace == nil {
		pc := make([]uintptr, 10)
		n := runtime.Callers(2, pc)
		if n > 0 {
			frames := runtime.CallersFrames(pc[:n])
			e.StackTrace = make([]string, 0, n)
			for {
				frame, more := frames.Next()
				e.StackTrace = append(e.StackTrace, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
				if !more {
					break
				}
			}
		}
	}
	return e
}

// WithService 标记服务名称
func (e *AppError) WithService(name string) *AppError {
	e.ServiceName = name
	return e
}

// New 创建新错误
func New(typ ErrorType, msg string, opts ...ErrorOption) *AppError {
	err := &AppError{
		Type:    typ,
		Code:    typeToCode(typ),
		Message: msg,
	}

	for _, opt := range opts {
		opt(err)
	}

	return err
}

// ErrorOption 错误配置选项
type ErrorOption func(*AppError)

// WithDetail -
func WithDetail(detail string) ErrorOption {
	return func(e *AppError) {
		e.Detail = detail
	}
}

// WithInternal -
func WithInternal(err error) ErrorOption {
	return func(e *AppError) {
		e.Internal = err
	}
}

// WithStack -
func WithStack() ErrorOption {
	return func(e *AppError) {
		e.WithStack()
	}
}

// 预定义错误构造器

// BadRequest -
func BadRequest(msg string, opts ...ErrorOption) *AppError {
	return New(ErrorTypeInvalidInput, msg, opts...)
}

// Unauthorized -
func Unauthorized(msg string, opts ...ErrorOption) *AppError {
	return New(ErrorTypeUnauthorized, msg, opts...)
}

// Forbidden -
func Forbidden(msg string, opts ...ErrorOption) *AppError {
	return New(ErrorTypeForbidden, msg, opts...)
}

// NotFound -
func NotFound(msg string, opts ...ErrorOption) *AppError {
	return New(ErrorTypeNotFound, msg, opts...)
}

// Conflict -
func Conflict(msg string, opts ...ErrorOption) *AppError {
	return New(ErrorTypeConflict, msg, opts...)
}

// RateLimited -
func RateLimited(msg string, opts ...ErrorOption) *AppError {
	return New(ErrorTypeRateLimited, msg, opts...)
}

// Internal -
func Internal(msg string, opts ...ErrorOption) *AppError {
	return New(ErrorTypeInternal, msg, opts...)
}

// NotImplemented -
func NotImplemented(msg string, opts ...ErrorOption) *AppError {
	return New(ErrorTypeNotImplemented, msg, opts...)
}

// Unavailable -
func Unavailable(msg string, opts ...ErrorOption) *AppError {
	return New(ErrorTypeUnavailable, msg, opts...)
}

// BindingFailed 新增绑定错误构造器
func BindingFailed(msg string, opts ...ErrorOption) *AppError {
	return New(ErrorTypeBindingFailed, msg, opts...)
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
	case ErrorTypeRateLimited:
		return http.StatusTooManyRequests
	case ErrorTypeInternal:
		return http.StatusInternalServerError
	case ErrorTypeNotImplemented:
		return http.StatusNotImplemented
	case ErrorTypeUnavailable:
		return http.StatusServiceUnavailable
	case ErrorTypeBindingFailed:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
