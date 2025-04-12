package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"io"
	"log/slog"
	"reflect"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	httperrors "Lumino/common/http_error_code"
	"github.com/gin-gonic/gin"
)

const LargeBody = "http: request body too large"

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
	isProduction := viper.GetString("env") == "production"

	// 如果是中间件的拦截，特殊处理一些部分
	if err.Internal != nil && strings.Contains(err.Internal.Error(), LargeBody) {
		err.Message = err.Message + LargeBody
	}

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

func handleBindError(err error) *httperrors.AppError {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var numError *strconv.NumError
	var verr validator.ValidationErrors

	switch {
	case errors.As(err, &syntaxError):
		return httperrors.BindingFailed("无效的JSON格式",
			httperrors.WithInternal(err),
			httperrors.WithDetail(fmt.Sprintf(
				"JSON语法错误 (位置: %d)",
				syntaxError.Offset,
			)),
		)

	case errors.As(err, &unmarshalTypeError):
		return formatTypeError(unmarshalTypeError)

	case errors.As(err, &numError):
		return httperrors.BindingFailed("数字格式错误",
			httperrors.WithInternal(err),
			httperrors.WithDetail(fmt.Sprintf(
				"字段值 '%s' 不是有效的数字",
				numError.Num,
			)),
		)

	case errors.Is(err, io.EOF):
		return httperrors.BindingFailed("请求体不能为空",
			httperrors.WithInternal(err),
			httperrors.WithDetail("empty_request_body"),
		)

	case errors.As(err, &verr):
		return formatValidationError(verr)

	default:
		return handleGenericTypeError(err)
	}
}

func formatTypeError(err *json.UnmarshalTypeError) *httperrors.AppError {
	return httperrors.BindingFailed("类型不匹配",
		httperrors.WithInternal(err),
		httperrors.WithDetail(fmt.Sprintf(
			"字段 '%s' 需要 %s 类型 (收到: %v)",
			err.Field,
			typeName(err.Type),
			err.Value,
		)),
	)
}

func typeName(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "整数"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "正整数"
	case reflect.Float32, reflect.Float64:
		return "小数"
	case reflect.Bool:
		return "布尔值"
	case reflect.String:
		return "字符串"
	case reflect.Struct:
		if t == reflect.TypeOf(time.Time{}) {
			return "日期时间"
		}
		return "对象"
	default:
		return t.String()
	}
}

func handleGenericTypeError(err error) *httperrors.AppError {
	// 尝试从错误信息中提取结构化信息
	if msg := extractTypeErrorMessage(err); msg != "" {
		return httperrors.BindingFailed("类型不匹配",
			httperrors.WithInternal(err),
			httperrors.WithDetail(msg),
		)
	}

	// 特殊处理时间解析错误
	if isTimeParseError(err) {
		return httperrors.BindingFailed("日期时间格式错误",
			httperrors.WithInternal(err),
			httperrors.WithDetail(cleanTimeErrorMessage(err.Error())),
		)
	}

	// 默认处理
	return httperrors.BindingFailed("无效的请求数据",
		httperrors.WithInternal(err),
		httperrors.WithDetail(cleanErrorMessage(err.Error())),
	)
}

func extractTypeErrorMessage(err error) string {
	// 处理 json.Unmarshal 类型错误模式
	re := regexp.MustCompile(`cannot unmarshal (.+?) into Go (struct field .+?) of type (.+)`)
	if matches := re.FindStringSubmatch(err.Error()); len(matches) == 4 {
		field := strings.TrimPrefix(matches[2], "struct field ")
		field = strings.Split(field, ".")[0] // 取最外层字段名
		return fmt.Sprintf("字段 '%s' 需要 %s 类型 (收到: %s)",
			field,
			parseTypeName(matches[3]),
			matches[1])
	}
	return ""
}

func parseTypeName(typeStr string) string {
	switch {
	case strings.Contains(typeStr, "int"):
		return "整数"
	case strings.Contains(typeStr, "float"):
		return "小数"
	case strings.Contains(typeStr, "bool"):
		return "布尔值"
	case strings.Contains(typeStr, "string"):
		return "字符串"
	case strings.Contains(typeStr, "time.Time"):
		return "日期时间"
	default:
		return typeStr
	}
}

func isTimeParseError(err error) bool {
	return strings.Contains(err.Error(), "time: ") ||
		strings.Contains(err.Error(), "parsing time ")
}

func cleanTimeErrorMessage(msg string) string {
	msg = strings.ReplaceAll(msg, "time: ", "")
	msg = strings.ReplaceAll(msg, "parsing time ", "")
	return strings.Trim(msg, `'"`)
}

func cleanErrorMessage(msg string) string {
	msg = strings.ReplaceAll(msg, "json: ", "")
	msg = strings.ReplaceAll(msg, "strconv: ", "")
	return msg
}

func formatValidationError(err validator.ValidationErrors) *httperrors.AppError {
	var details []string
	for _, e := range err {
		details = append(details, fmt.Sprintf(
			"字段 %s 验证失败 (%s)",
			e.Field(),
			e.Tag(),
		))
	}
	return httperrors.BindingFailed("请求数据验证失败: "+strings.Join(details, "; "),
		httperrors.WithInternal(err),
		httperrors.WithDetail(strings.Join(details, "; ")),
	)
}

// Bind 封装了ShouldBind并返回自定义错误
func Bind(c *gin.Context, obj any) error {
	if err := c.ShouldBind(obj); err != nil {
		return handleBindError(err)
	}
	return nil
}

// BindJSON 然后在BindJSON等函数中使用：
func BindJSON(c *gin.Context, obj any) error {
	return Bind(c, obj)
}

// BindQuery 封装了ShouldBindQuery
func BindQuery(c *gin.Context, obj any) error {
	return Bind(c, obj)
}

// BindURI 封装了ShouldBindUri
func BindURI(c *gin.Context, obj any) error {
	return Bind(c, obj)
}
