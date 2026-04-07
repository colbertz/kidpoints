package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"points-app/service"
)

type DrawHandler struct {
	drawService *service.DrawService
}

func NewDrawHandler(drawService *service.DrawService) *DrawHandler {
	return &DrawHandler{drawService: drawService}
}

type DrawRequest struct {
	KidID int `json:"kid_id"`
}

func (h *DrawHandler) Draw(c *gin.Context) {
	var req DrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	prize, err := h.drawService.Draw(req.KidID)
	if err != nil {
		if errors.Is(err, service.ErrKidNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "kid not found"})
			return
		}
		if errors.Is(err, service.ErrInsufficientPointsForDraw) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient points for draw"})
			return
		}
		if errors.Is(err, service.ErrNoPrizesAvailable) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no prizes available"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prize)
}
