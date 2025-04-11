package middleware

import (
	"Lumino/common/http_error_code"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

// HttpErrorHandler Gin错误处理中间件
func HttpErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 先处理请求

		// 处理显式添加的错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			appErr := toAppError(err)

			// 记录内部错误日志
			if appErr.Internal != nil {
				log.Printf("Error: %v\nStack: %+v", appErr.Internal, appErr)
			}

			c.JSON(appErr.Code, gin.H{
				"type":    appErr.Type,
				"message": appErr.Message,
				"detail":  appErr.Detail,
			})
			return
		}

		// 处理panic恢复
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch v := r.(type) {
				case error:
					err = v
				case string:
					err = fmt.Errorf(v)
				default:
					err = fmt.Errorf("%v", v)
				}

				appErr := toAppError(err)
				log.Printf("Panic recovered: %+v", err)

				c.JSON(appErr.Code, gin.H{
					"type":    appErr.Type,
					"message": appErr.Message,
					"detail":  "internal server error",
				})
			}
		}()
	}
}

// toAppError 转换各种错误为AppError
func toAppError(err error) *http_error_code.AppError {
	switch e := err.(type) {
	case *http_error_code.AppError:
		return e
	default:
		// 处理GORM错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http_error_code.NotFound("资源不存在", "", err)
		}

		// 默认作为服务器内部错误
		return http_error_code.InternalServer("服务器内部错误", "", err)
	}
}
