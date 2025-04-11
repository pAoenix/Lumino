package middleware

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"
	"time"

	httperrors "Lumino/common/http_error_code"
	"github.com/gin-gonic/gin"
)

// ErrorHandler 通用错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 处理panic恢复
		defer func() {
			if r := recover(); r != nil {
				// 记录堆栈
				stack := debug.Stack()

				// 构建错误
				err := buildRecoveryError(r, stack)

				// 响应客户端
				respondError(c, err)
			}
		}()

		c.Next() // 处理请求

		// 处理请求错误
		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last().Err
			appErr := toAppError(lastErr)

			// 响应客户端
			respondError(c, appErr)
		}
	}
}

// buildRecoveryError 从panic构建错误
func buildRecoveryError(r any, stack []byte) *httperrors.AppError {
	var err error
	switch v := r.(type) {
	case error:
		err = v
	case string:
		err = fmt.Errorf(v)
	default:
		err = fmt.Errorf("%v", v)
	}

	return httperrors.Internal("服务器内部错误",
		httperrors.WithInternal(err),
		httperrors.WithDetail(string(stack)),
	)
}

// toAppError 转换各种错误为AppError
func toAppError(err error) *httperrors.AppError {
	switch e := err.(type) {
	case *httperrors.AppError:
		return e
	case validator.ValidationErrors:
		// 处理验证错误
		return formatValidationError(e)
	default:
		// 可以添加更多特定错误的转换逻辑
		return httperrors.Internal("服务器内部错误",
			httperrors.WithInternal(e),
		)
	}
}

// logError 记录错误日志
func logError(logger *slog.Logger, c *gin.Context, err *httperrors.AppError, start time.Time) {
	duration := time.Since(start)

	attrs := []slog.Attr{
		slog.String("method", c.Request.Method),
		slog.String("path", c.Request.URL.Path),
		slog.String("ip", c.ClientIP()),
		slog.Duration("duration", duration),
		slog.String("error_type", string(err.Type)),
		slog.Int("status_code", err.Code),
	}

	if err.Internal != nil {
		attrs = append(attrs, slog.String("internal_error", err.Internal.Error()))
	}

	if len(err.StackTrace) > 0 {
		attrs = append(attrs, slog.Any("stack", err.StackTrace))
	}

	logger.LogAttrs(
		context.Background(),
		errorLevel(err.Code),
		err.Message,
		attrs...,
	)
}

// errorLevel 根据状态码确定日志级别
func errorLevel(code int) slog.Level {
	switch {
	case code >= 500:
		return slog.LevelError
	case code >= 400:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

// respondError 响应错误给客户端
func respondError(c *gin.Context, err *httperrors.AppError) {
	// 在生产环境隐藏内部细节
	isProduction := os.Getenv("ENV") == "production"

	response := gin.H{
		"type":    err.Type,
		"message": err.Message,
	}

	// 非生产环境返回更多调试信息
	if !isProduction {
		if err.Detail != "" {
			response["detail"] = err.Detail
		}
	}

	c.JSON(err.Code, response)
}

// RegisterCustomValidators 注册自定义验证器
func RegisterCustomValidators(v *validator.Validate) {
	// 示例：注册自定义验证器
	v.RegisterValidation("strong_password", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) >= 8
	})
}

// formatValidationError 格式化验证错误
func formatValidationError(err validator.ValidationErrors) *httperrors.AppError {
	var details []string
	for _, e := range err {
		details = append(details, fmt.Sprintf(
			"字段 %s 验证失败 (%s)",
			e.Field(),
			e.Tag(),
		))
	}
	return httperrors.BindingFailed("请求数据验证失败",
		httperrors.WithInternal(err),
		httperrors.WithDetail(strings.Join(details, "; ")),
	)
}

// BindJSON 然后在BindJSON等函数中使用：
func BindJSON(c *gin.Context, obj any) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			return formatValidationError(verr)
		}
		return httperrors.BindingFailed("无效的请求数据",
			httperrors.WithInternal(err),
			httperrors.WithDetail("json_binding_error"),
		)
	}
	return nil
}

// Bind 封装了ShouldBind并返回自定义错误
func Bind(c *gin.Context, obj any) error {
	if err := c.ShouldBind(obj); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			return formatValidationError(verr)
		}
		return httperrors.BindingFailed("无效的参数数据",
			httperrors.WithInternal(err),
			httperrors.WithDetail("json_binding_error"),
		)
	}
	return nil
}

// BindQuery 封装了ShouldBindQuery
func BindQuery(c *gin.Context, obj any) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			return formatValidationError(verr)
		}
		return httperrors.BindingFailed("无效的查询数据",
			httperrors.WithInternal(err),
			httperrors.WithDetail("json_binding_error"),
		)
	}
	return nil
}

// BindURI 封装了ShouldBindUri
func BindURI(c *gin.Context, obj any) error {
	if err := c.ShouldBindUri(obj); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			return formatValidationError(verr)
		}
		return httperrors.BindingFailed("无效的URI数据",
			httperrors.WithInternal(err),
			httperrors.WithDetail("json_binding_error"),
		)
	}
	return nil
}
