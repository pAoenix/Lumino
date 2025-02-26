package router

import (
	"Lumino/common/logger"
	"Lumino/router/middleware"
	"Lumino/server"
	"Lumino/store"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"net/http"
)

// Router -
type Router struct {
	fx.In
	DB                *store.DB
	HealthServer      *server.HealthServer
	TransactionServer *server.TransactionServer
	UserServer        *server.UserServer
	CategoryServer    *server.CategoryServer
	AccountBookServer *server.AccountBookServer
}

// Handler -
func (r *Router) Handler() http.Handler {
	gin.DisableConsoleColor()
	e := gin.New()
	e.Use(middleware.Cors())
	e.Use(gin.Recovery())
	e.Use(middleware.Log(logger.Logger))
	e.Use(middleware.DB(r.DB))
	e.GET("api/v1/health", r.HealthServer.Health)

	e.GET("api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	transaction := e.Group("api/v1/transaction")
	{
		transaction.POST("", r.TransactionServer.Register)
		transaction.GET("", r.TransactionServer.Get)
	}
	user := e.Group("api/v1/user")
	{
		user.POST("", r.UserServer.Register)
	}

	category := e.Group("api/v1/category")
	{
		category.POST("", r.CategoryServer.Register)
		category.GET("", r.CategoryServer.Get)
	}

	accountBook := e.Group("api/v1/account-book")
	{
		accountBook.GET("", r.AccountBookServer.Get)
	}

	return e
}

// NewHttpServer -
func NewHttpServer(router Router) *http.Server {
	return &http.Server{
		Addr:    ":" + viper.GetString("port"),
		Handler: router.Handler(),
	}
}
