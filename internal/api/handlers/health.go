package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck returns the health status of the service
// @Summary Health check
// @Description Returns the health status of the blind service
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "blind-service",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// ReadinessCheck returns the readiness status of the service
// @Summary Readiness check
// @Description Returns whether the service is ready to accept requests
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /ready [get]
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// LivenessCheck returns the liveness status of the service
// @Summary Liveness check
// @Description Returns whether the service is alive
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /live [get]
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}
