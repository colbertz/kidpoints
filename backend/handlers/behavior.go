package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"points-app/models"
	"points-app/repository"
)

type BehaviorHandler struct {
	repo *repository.BehaviorRepository
}

func NewBehaviorHandler(repo *repository.BehaviorRepository) *BehaviorHandler {
	return &BehaviorHandler{repo: repo}
}

func (h *BehaviorHandler) GetBehaviors(c *gin.Context) {
	behaviors, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, behaviors)
}

func (h *BehaviorHandler) CreateBehavior(c *gin.Context) {
	var b models.Behavior
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.repo.Create(&b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, b)
}

func (h *BehaviorHandler) UpdateBehavior(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var b models.Behavior
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	b.ID = id

	if err := h.repo.Update(&b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *BehaviorHandler) DeleteBehavior(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
