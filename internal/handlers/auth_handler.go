package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/radhikadarode/finance-backend/internal/models"
	"github.com/radhikadarode/finance-backend/internal/services"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{service: services.NewAuthService()}
}

// Login godoc
// @Summary Login and receive a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "Credentials"
// @Success 200 {object} models.APIResponse
// @Failure 400,401 {object} models.APIResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	resp, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data:    resp,
	})
}
