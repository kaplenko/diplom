package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kaplenko/diplom/internal/models"
	"github.com/kaplenko/diplom/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetMe godoc
// @Summary      Get current user
// @Description  Get the profile of the authenticated user
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.UserResponse
// @Failure      401  {object}  models.ErrorResponse
// @Router       /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := getUserID(c)

	user, err := h.userService.GetByID(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// UpdateMe godoc
// @Summary      Update current user
// @Description  Update the profile of the authenticated user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      models.UpdateUserRequest  true  "Update data"
// @Success      200   {object}  models.UserResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      401   {object}  models.ErrorResponse
// @Router       /users/me [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := getUserID(c)

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err.Error())
		return
	}

	user, err := h.userService.Update(userID, req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// ListUsers godoc
// @Summary      List users (admin)
// @Description  Get a paginated list of all users
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Param        page       query  int     false  "Page number"  default(1)
// @Param        page_size  query  int     false  "Page size"    default(10)
// @Param        search     query  string  false  "Search by name or email"
// @Success      200  {object}  models.PaginatedResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Router       /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	params := parsePagination(c)

	users, total, err := h.userService.List(params)
	if err != nil {
		respondError(c, err)
		return
	}

	responses := make([]models.UserResponse, len(users))
	for i, u := range users {
		responses[i] = u.ToResponse()
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:     responses,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	})
}

// DeleteUser godoc
// @Summary      Delete user (admin)
// @Description  Delete a user by ID
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "User ID"
// @Success      204
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid user ID")
		return
	}

	if err := h.userService.Delete(id); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
