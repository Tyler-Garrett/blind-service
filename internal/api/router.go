package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tyler-garrett/blind-service/internal/api/handlers"
	"github.com/tyler-garrett/blind-service/internal/api/middleware"
	"github.com/tyler-garrett/blind-service/internal/config"
)

func SetupRouter(
	cfg *config.Config,
	healthHandler *handlers.HealthHandler,
	blindHandler *handlers.BlindHandler,
) *gin.Engine {
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware(&cfg.CORS))

	// no /api prefix for k8s
	router.GET("/health", healthHandler.HealthCheck)
	router.GET("/ready", healthHandler.ReadinessCheck)
	router.GET("/live", healthHandler.LivenessCheck)

	v1 := router.Group("/api/v1")
	{
		blinds := v1.Group("/blinds")
		{
			blinds.GET("", blindHandler.GetAll)
			blinds.GET("/stats", blindHandler.GetStats)
			blinds.GET("/:id", blindHandler.GetById)
			blinds.POST("", blindHandler.Create)
			blinds.POST("/radius", blindHandler.GetInRadius)
			blinds.PUT("/:id", blindHandler.Update)
			blinds.DELETE("/:id", blindHandler.Delete)
		}
	}

	return router
}
