package router

import (
	"feed_system/internal/bootstrap"
	"feed_system/internal/handler"
	"feed_system/internal/middleware"
	"feed_system/internal/response"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New(app *bootstrap.App) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.ErrorHandler())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.DeepLinking(true),
		ginSwagger.DocExpansion("none"),
		ginSwagger.DefaultModelsExpandDepth(-1),
	))

	// Health godoc
	// @Summary Health check
	// @Description Check whether the HTTP server is running.
	// @Tags system
	// @Produce json
	// @Success 200 {object} response.Body
	// @Router /health [get]
	r.GET("/health", middleware.Wrap(func(c *gin.Context) error {
		response.Success(c, gin.H{
			"status": "ok",
		})
		return nil
	}))

	authHandler := handler.NewAuthHandler(app.Auth)
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
		}
	}

	return r
}
