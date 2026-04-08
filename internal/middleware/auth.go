package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/kaplenko/diplom/internal/models"
	"github.com/kaplenko/diplom/internal/service"
)

func AuthRequired(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: models.ErrorBody{Code: "UNAUTHORIZED", Message: "Authorization header is required"},
			})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: models.ErrorBody{Code: "UNAUTHORIZED", Message: "Authorization header must be: Bearer <token>"},
			})
			return
		}

		userID, role, err := authService.ParseAccessToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: models.ErrorBody{Code: "UNAUTHORIZED", Message: "Invalid or expired token"},
			})
			return
		}

		c.Set("user_id", userID)
		c.Set("user_role", role)
		c.Next()
	}
}
