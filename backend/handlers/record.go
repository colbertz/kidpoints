package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"points-app/service"
)

type RecordHandler struct {
	recordService *service.RecordService
}

func NewRecordHandler(recordService *service.RecordService) *RecordHandler {
	return &RecordHandler{recordService: recordService}
}

func (h *RecordHandler) GetRecords(c *gin.Context) {
	kidIDStr := c.Query("kid_id")
	if kidIDStr != "" {
		kidID, err := strconv.Atoi(kidIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kid_id"})
			return
		}
		records, err := h.recordService.GetRecordsByKid(kidID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, records)
		return
	}

	records, err := h.recordService.GetAllRecords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}

type ReverseRequest struct {
	Reason string `json:"reason"`
}

func (h *RecordHandler) ReverseRecord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record id"})
		return
	}

	var req ReverseRequest
	c.ShouldBindJSON(&req)

	record, err := h.recordService.ReverseRecord(id, req.Reason)
	if err != nil {
		if errors.Is(err, service.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}
		if errors.Is(err, service.ErrRecordAlreadyReversed) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "record already reversed"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, record)
}
