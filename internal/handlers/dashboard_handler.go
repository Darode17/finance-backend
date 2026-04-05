package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/radhikadarode/finance-backend/internal/models"
	"github.com/radhikadarode/finance-backend/internal/services"
)

type DashboardHandler struct {
	service *services.DashboardService
}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{service: services.NewDashboardService()}
}

// GetSummary - Analyst and Admin
func (h *DashboardHandler) GetSummary(c *gin.Context) {
	summary, err := h.service.GetSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: summary})
}
