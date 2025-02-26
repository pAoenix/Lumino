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
