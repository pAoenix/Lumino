package common

import "github.com/gin-gonic/gin"

const (
	DebugMode   = gin.DebugMode
	ReleaseMode = gin.ReleaseMode
	TestMode    = gin.TestMode
)

// Mode 服务运行模式/debug/release/test
func Mode() string {
	return gin.Mode()
}
