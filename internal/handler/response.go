package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kaplenko/diplom/internal/models"
	"github.com/kaplenko/diplom/internal/service"
)

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: models.ErrorBody{Code: "NOT_FOUND", Message: "Resource not found"},
		})
	case errors.Is(err, service.ErrConflict):
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error: models.ErrorBody{Code: "CONFLICT", Message: "Resource already exists"},
		})
	case errors.Is(err, service.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorBody{Code: "UNAUTHORIZED", Message: "Invalid credentials"},
		})
	case errors.Is(err, service.ErrForbidden):
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Error: models.ErrorBody{Code: "FORBIDDEN", Message: "Access denied"},
		})
	case errors.Is(err, service.ErrInvalidToken):
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorBody{Code: "UNAUTHORIZED", Message: "Invalid or expired token"},
		})
	default:
		log.Printf("internal error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorBody{Code: "INTERNAL_ERROR", Message: "Internal server error"},
		})
	}
}

func respondValidationError(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, models.ErrorResponse{
		Error: models.ErrorBody{Code: "VALIDATION_ERROR", Message: msg},
	})
}

func parsePagination(c *gin.Context) models.PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")

	p := models.PaginationParams{Page: page, PageSize: pageSize, Search: search}
	p.Normalize()
	return p
}

func getUserID(c *gin.Context) int64 {
	id, ok := c.Get("user_id")
	if !ok {
		return 0
	}
	uid, ok := id.(int64)
	if !ok {
		return 0
	}
	return uid
}
