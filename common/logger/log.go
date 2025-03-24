package logger

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var zapLogger *zap.SugaredLogger

// init 代码运行提前初始化，自动运行
func init() {
	// 从配置危机
	logLevel := viper.GetString("log.level")
	// 设置日志级别
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		panic(any(err))
	}
	// 配置日志
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	Logger, err = cfg.Build()
	if err != nil {
		panic(any(err))
	}
	zapLogger = Logger.Sugar()
}

// Error -
func Error(args ...interface{}) {
	zapLogger.Error(args)
}

// Info -
func Info(args ...interface{}) {
	zapLogger.Info(args)
}

// Warn -
func Warn(args ...interface{}) {
	zapLogger.Warn(args)
}

// Fatalf -
func Fatalf(msg string, args ...interface{}) {
	zapLogger.Fatalf(msg, args)
}
