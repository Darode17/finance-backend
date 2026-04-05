package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/radhikadarode/finance-backend/internal/models"
	"github.com/radhikadarode/finance-backend/internal/services"
)

type RecordHandler struct {
	service *services.RecordService
}

func NewRecordHandler() *RecordHandler {
	return &RecordHandler{service: services.NewRecordService()}
}

// CreateRecord - Admin and Analyst
func (h *RecordHandler) CreateRecord(c *gin.Context) {
	var req models.CreateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	record, err := h.service.CreateRecord(req, c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Record created successfully",
		Data:    record,
	})
}

// GetRecords - All authenticated users
func (h *RecordHandler) GetRecords(c *gin.Context) {
	var filter models.RecordFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	result, err := h.service.GetRecords(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: result})
}

// GetRecord - All authenticated users
func (h *RecordHandler) GetRecord(c *gin.Context) {
	record, err := h.service.GetRecordByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{Success: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: record})
}

// UpdateRecord - Admin only
func (h *RecordHandler) UpdateRecord(c *gin.Context) {
	var req models.UpdateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	record, err := h.service.UpdateRecord(c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Record updated successfully",
		Data:    record,
	})
}

// DeleteRecord - Admin only
func (h *RecordHandler) DeleteRecord(c *gin.Context) {
	if err := h.service.DeleteRecord(c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{Success: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "Record deleted successfully"})
}
