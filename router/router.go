package router

import (
	"Lumino/common/logger"
	"Lumino/docs"
	"Lumino/router/middleware"
	"Lumino/server"
	"Lumino/store"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"net/http"
	"strings"
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
	FriendServer      *server.FriendServer
	AccountServer     *server.AccountServer
}

// Handler -
func (r *Router) Handler() http.Handler {
	setupValidator()
	gin.DisableConsoleColor()
	e := gin.New()
	e.Use(middleware.SizeLimitMiddleware())
	e.Use(middleware.DB(r.DB))
	e.Use(middleware.Cors())
	e.Use(middleware.Log(logger.Logger))
	e.Use(middleware.ErrorHandler())
	e.Use(gin.Recovery())
	e.GET("api/v1/health", r.HealthServer.Health)
	docs.SwaggerInfo.BasePath = ""
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	transaction := e.Group("api/v1/transaction")
	{
		transaction.POST("", r.TransactionServer.Register)
		transaction.GET("", r.TransactionServer.Get)
		transaction.PUT("", r.TransactionServer.Modify)
		transaction.DELETE("", r.TransactionServer.Delete)
	}
	user := e.Group("api/v1/user")
	{
		user.POST("", r.UserServer.Register)
		user.PUT("", r.UserServer.Modify)
		user.PUT("profile-photo", r.UserServer.ModifyProfilePhoto)
		user.GET("", r.UserServer.Get)
		user.DELETE("", r.UserServer.Delete)
	}

	friend := e.Group("api/v1/friend")
	{
		friend.POST("/invite", r.FriendServer.Invite)
		friend.DELETE("", r.FriendServer.Delete)
	}

	category := e.Group("api/v1/category")
	{
		category.POST("", r.CategoryServer.Register)
		category.GET("", r.CategoryServer.Get)
		category.PUT("", r.CategoryServer.Modify)
		category.PUT("/icon-image", r.CategoryServer.ModifyIconImage)
		category.DELETE("", r.CategoryServer.Delete)
	}

	accountBook := e.Group("api/v1/account-book")
	{
		accountBook.GET(":id", r.AccountBookServer.GetByID)
		accountBook.GET("", r.AccountBookServer.Get)
		accountBook.POST("", r.AccountBookServer.Register)
		accountBook.POST("/merge", r.AccountBookServer.Merge)
		accountBook.PUT("", r.AccountBookServer.Modify)
		accountBook.DELETE("", r.AccountBookServer.Delete)
	}

	account := e.Group("/api/v1/account")
	{
		account.POST("", r.AccountServer.Register)
		account.PUT("", r.AccountServer.Modify)
		account.GET("", r.AccountServer.Get)
		account.DELETE("", r.AccountServer.Delete)
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

func setupValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册 notblank 验证规则
		_ = v.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			return len(strings.TrimSpace(value)) > 0
		})
	}
}
