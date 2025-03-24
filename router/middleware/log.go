package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// Log gon-log format
func Log(logger *zap.Logger, ignorePathList ...string) gin.HandlerFunc {
	ignores := make(map[string]struct{})
	for _, ignorePath := range ignorePathList {
		ignores[ignorePath] = struct{}{}
	}
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		if _, ok := ignores[path]; ok {
			return
		}
		end := time.Now()
		latency := time.Since(start)
		fields := []zap.Field{
			zap.String("time", end.Format(time.RFC3339)),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		}
		logger.Info(path, fields...)
	}
}
