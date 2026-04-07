package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"points-app/models"
	"points-app/service"
)

type KidHandler struct {
	kidService *service.KidService
}

func NewKidHandler(kidService *service.KidService) *KidHandler {
	return &KidHandler{kidService: kidService}
}

func (h *KidHandler) GetKids(c *gin.Context) {
	kids, err := h.kidService.GetAllKids()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kids)
}

func (h *KidHandler) GetPoints(c *gin.Context) {
	nameOrID := c.Param("id")
	kidID, err := h.resolveKidID(nameOrID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "kid not found"})
		return
	}

	points, err := h.kidService.GetKidPoints(kidID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"points": points})
}

type PointsRequest struct {
	ID         int    `json:"id"`
	Name      string `json:"name"`
	BehaviorID int    `json:"behavior_id"`
}

func (h *KidHandler) AddPoints(c *gin.Context) {
	var req PointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	kidID, err := h.resolveRequestKidID(req.ID, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kid id or name required"})
		return
	}

	if err := h.kidService.AddPoints(kidID, req.BehaviorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	kid, _ := h.kidService.GetKid(kidID)
	c.JSON(http.StatusOK, kid)
}

func (h *KidHandler) SubPoints(c *gin.Context) {
	var req PointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	kidID, err := h.resolveRequestKidID(req.ID, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kid id or name required"})
		return
	}

	if err := h.kidService.SubPoints(kidID, req.BehaviorID); err != nil {
		if errors.Is(err, service.ErrInsufficientPoints) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient points"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	kid, _ := h.kidService.GetKid(kidID)
	c.JSON(http.StatusOK, kid)
}

func (h *KidHandler) resolveKidID(nameOrID string) (int, error) {
	if id, err := parseInt(nameOrID); err == nil {
		return id, nil
	}
	kid, err := h.kidService.GetKidByName(nameOrID)
	if err != nil {
		return 0, err
	}
	return kid.ID, nil
}

func (h *KidHandler) resolveRequestKidID(id int, name string) (int, error) {
	if id > 0 {
		return id, nil
	}
	if name != "" {
		kid, err := h.kidService.GetKidByName(name)
		if err != nil {
			return 0, err
		}
		return kid.ID, nil
	}
	return 0, errors.New("kid id or name required")
}

func parseInt(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, errors.New("not a number")
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

func (h *KidHandler) CreateKid(c *gin.Context) {
	var kid models.Kid
	if err := c.ShouldBindJSON(&kid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if kid.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if kid.Avatar == "" {
		kid.Avatar = "👤"
	}
	if kid.Points < 0 {
		kid.Points = 0
	}
	if err := h.kidService.CreateKid(&kid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, kid)
}

func (h *KidHandler) UpdateKid(c *gin.Context) {
	idStr := c.Param("id")
	id, err := parseInt(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kid id"})
		return
	}

	var kid models.Kid
	if err := c.ShouldBindJSON(&kid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	kid.ID = id

	if err := h.kidService.UpdateKid(&kid); err != nil {
		if errors.Is(err, service.ErrKidNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "kid not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kid)
}

func (h *KidHandler) DeleteKid(c *gin.Context) {
	idStr := c.Param("id")
	id, err := parseInt(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kid id"})
		return
	}

	if err := h.kidService.DeleteKid(id); err != nil {
		if errors.Is(err, service.ErrKidNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "kid not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
