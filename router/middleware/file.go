package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// SizeLimitMiddleware -
func SizeLimitMiddleware(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在请求进入时立即限制
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
