package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kaplenko/diplom/internal/models"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		r, ok := role.(models.Role)
		if !exists || !ok || r != models.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse{
				Error: models.ErrorBody{Code: "FORBIDDEN", Message: "Admin access required"},
			})
			return
		}
		c.Next()
	}
}
