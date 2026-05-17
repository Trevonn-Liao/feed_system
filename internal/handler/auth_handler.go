package handler

import (
	"feed_system/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Register godoc
// @Summary Account password register
// @Description Register by username and password only.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body service.RegisterRequest true "register payload"
// @Success 200 {object} response.Body
// @Failure 400 {object} response.Body
// @Failure 409 {object} response.Body
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	if err := h.auth.Register(c); err != nil {
		_ = c.Error(err)
	}
}

// Login godoc
// @Summary Account password login
// @Description Login by username and password only.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body service.LoginRequest true "login payload"
// @Success 200 {object} response.Body
// @Failure 400 {object} response.Body
// @Failure 401 {object} response.Body
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	if err := h.auth.Login(c); err != nil {
		_ = c.Error(err)
	}
}

// Refresh godoc
// @Summary Refresh token pair
// @Description Use refresh token to get a new token pair.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body service.RefreshRequest true "refresh payload"
// @Success 200 {object} response.Body
// @Failure 400 {object} response.Body
// @Failure 401 {object} response.Body
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	if err := h.auth.Refresh(c); err != nil {
		_ = c.Error(err)
	}
}
