package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

// SizeLimitMiddleware -
func SizeLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在请求进入时立即限制
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, viper.GetInt64("bodySizeLimit"))
		c.Next()
	}
}
