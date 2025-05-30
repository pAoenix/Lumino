package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

// HealthServer -
type HealthServer struct {
}

// NewHealthServer -
func NewHealthServer() *HealthServer {
	return &HealthServer{}
}

// Health 健康检查
// @Summary	健康检查
// @Tags 健康检查
// @Success	200
// @Failure	500 {object}  http_error_code.AppError      "服务端异常"
// @Router		/api/v1/health [get]
func (s *HealthServer) Health(c *gin.Context) {
	hostname, _ := os.Hostname()
	c.JSON(http.StatusOK, gin.H{
		"service_name": "Lumino",
		"hostname":     hostname,
		"message":      "success",
		"env":          viper.GetString("env"),
	})
	return
}
