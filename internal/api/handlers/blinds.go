package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tyler-garrett/blind-service/internal/domain"
	"github.com/tyler-garrett/blind-service/internal/service"
)

type BlindHandler struct {
	service *service.BlindService
}

func NewBlindHandler(service *service.BlindService) *BlindHandler {
	return &BlindHandler{
		service: service,
	}
}

// @Summary Get all blinds
// @Description Retrieves all registered waterfowl blinds in Virginia
// @Tags blinds
// @Produce json
// @Success 200 {array} domain.Blind
// @Failure 500 {object} map[string]string
// @Router /blinds [get]
func (h *BlindHandler) GetAll(c *gin.Context) {
	blinds, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blinds)
}

// @Summary Get blind by Id
// @Description Retrieves a specific blind by its Id
// @Tags blinds
// @Produce json
// @Param id path string true "Blind Id"
// @Success 200 {object} domain.Blind
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /blinds/{id} [get]
func (h *BlindHandler) GetById(c *gin.Context) {
	id := c.Param("id")

	blind, err := h.service.GetById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrBlindNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "blind not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blind)
}

// GetInRadius retrieves blinds within a radius
// @Summary Get blinds within radius
// @Description Retrieves all blinds within a specified radius of a location
// @Tags blinds
// @Accept json
// @Produce json
// @Param request body domain.BlindsInRadiusRequest true "Location and radius"
// @Success 200 {array} domain.Blind
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /blinds/radius [post]
func (h *BlindHandler) GetInRadius(c *gin.Context) {
	var req domain.BlindsInRadiusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blinds, err := h.service.GetInRadius(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blinds)
}

// Create creates a new blind
// @Summary Create a new blind
// @Description Creates a new waterfowl blind registration
// @Tags blinds
// @Accept json
// @Produce json
// @Param request body domain.CreateBlindRequest true "Blind data"
// @Success 201 {object} domain.Blind
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /blinds [post]
func (h *BlindHandler) Create(c *gin.Context) {
	var req domain.CreateBlindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blind, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrBlindAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "blind already exists"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, blind)
}

// Update updates a blind
// @Summary Update a blind
// @Description Updates an existing blind's location
// @Tags blinds
// @Accept json
// @Produce json
// @Param id path string true "Blind Id"
// @Param request body domain.UpdateBlindRequest true "Updated blind data"
// @Success 200 {object} domain.Blind
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /blinds/{id} [put]
func (h *BlindHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req domain.UpdateBlindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blind, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, domain.ErrBlindNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "blind not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blind)
}

// Delete deletes a blind
// @Summary Delete a blind
// @Description Deletes a blind registration
// @Tags blinds
// @Produce json
// @Param id path string true "Blind Id"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /blinds/{id} [delete]
func (h *BlindHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrBlindNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "blind not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetStats retrieves statistics about blinds
// @Summary Get blind statistics
// @Description Retrieves statistical information about all blinds
// @Tags blinds
// @Produce json
// @Success 200 {object} domain.BlindStats
// @Failure 500 {object} map[string]string
// @Router /blinds/stats [get]
func (h *BlindHandler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
} //hi
