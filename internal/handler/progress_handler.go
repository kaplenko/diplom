package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kaplenko/diplom/internal/service"
)

type ProgressHandler struct {
	progressService *service.ProgressService
}

func NewProgressHandler(progressService *service.ProgressService) *ProgressHandler {
	return &ProgressHandler{progressService: progressService}
}

// GetCourseProgress godoc
// @Summary      Get course progress
// @Description  Get the authenticated user's progress for a specific course
// @Tags         Progress
// @Produce      json
// @Security     BearerAuth
// @Param        course_id  path  int  true  "Course ID"
// @Success      200  {object}  models.CourseProgress
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /courses/{course_id}/progress [get]
func (h *ProgressHandler) GetCourseProgress(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("course_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid course ID")
		return
	}

	cp, err := h.progressService.GetCourseProgress(getUserID(c), courseID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, cp)
}

// GetAllProgress godoc
// @Summary      Get all progress
// @Description  Get the authenticated user's progress across all courses
// @Tags         Progress
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  models.CourseProgress
// @Failure      401  {object}  models.ErrorResponse
// @Router       /progress [get]
func (h *ProgressHandler) GetAllProgress(c *gin.Context) {
	list, err := h.progressService.GetAllProgress(getUserID(c))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, list)
}
