package router

import (
	"feed_system/internal/middleware"
	"feed_system/internal/response"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New() *gin.Engine {
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

	return r
}
