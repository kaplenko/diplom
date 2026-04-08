package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kaplenko/diplom/internal/models"
	"github.com/kaplenko/diplom/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new student account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.RegisterRequest  true  "Registration data"
// @Success      201   {object}  models.UserResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      409   {object}  models.ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err.Error())
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user.ToResponse())
}

// Login godoc
// @Summary      Login
// @Description  Authenticate user and return JWT tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.LoginRequest  true  "Login credentials"
// @Success      200   {object}  models.TokenResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      401   {object}  models.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err.Error())
		return
	}

	tokens, err := h.authService.Login(req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Refresh godoc
// @Summary      Refresh tokens
// @Description  Get new access and refresh tokens using a valid refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.RefreshRequest  true  "Refresh token"
// @Success      200   {object}  models.TokenResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      401   {object}  models.ErrorResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err.Error())
		return
	}

	tokens, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}
