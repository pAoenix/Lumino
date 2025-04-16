package main

import (
	"Lumino/common"
	"Lumino/common/logger"
	"Lumino/config"
	"Lumino/router"
	"Lumino/server"
	"Lumino/service"
	"Lumino/store"
	"context"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"net/http"
)

// HttpServerLifetimeHook 设置生命周期钩子函数
func HttpServerLifetimeHook(lc fx.Lifecycle, srv *http.Server) {
	// 设置生命周期钩子函数
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Error("fail to listen port", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("HTTP服务器停止中...")
			ctx, cancel := context.WithTimeout(ctx, viper.GetDuration("gracefulShutdown.timeoutSeconds"))
			defer cancel()
			return srv.Shutdown(ctx)
		},
	})
}

// @title						Lumino
// @version					1.0
// @description				lumino 记账本服务
// @host						localhost:8080
// @BasePath					/api/v1
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	config.LoadConfig()
	app := fx.New(
		common.Module,
		store.Module,
		service.Module,
		server.Module,
		router.Module,
		fx.Invoke(HttpServerLifetimeHook),
	)
	app.Run()
}
