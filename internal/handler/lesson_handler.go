package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kaplenko/diplom/internal/models"
	"github.com/kaplenko/diplom/internal/service"
)

type LessonHandler struct {
	lessonService *service.LessonService
}

func NewLessonHandler(lessonService *service.LessonService) *LessonHandler {
	return &LessonHandler{lessonService: lessonService}
}

// CreateLesson godoc
// @Summary      Create a lesson (admin)
// @Description  Create a new lesson in a course
// @Tags         Lessons
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        course_id  path  int                          true  "Course ID"
// @Param        body       body  models.CreateLessonRequest   true  "Lesson data"
// @Success      201  {object}  models.Lesson
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /courses/{course_id}/lessons [post]
func (h *LessonHandler) Create(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("course_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid course ID")
		return
	}

	var req models.CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err.Error())
		return
	}

	lesson, err := h.lessonService.Create(courseID, req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, lesson)
}

// GetLesson godoc
// @Summary      Get a lesson
// @Description  Get lesson by ID
// @Tags         Lessons
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Lesson ID"
// @Success      200  {object}  models.Lesson
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /lessons/{id} [get]
func (h *LessonHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("lesson_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid lesson ID")
		return
	}

	lesson, err := h.lessonService.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, lesson)
}

// ListLessons godoc
// @Summary      List lessons in a course
// @Description  Get a paginated list of lessons for a course
// @Tags         Lessons
// @Produce      json
// @Security     BearerAuth
// @Param        course_id  path   int     true   "Course ID"
// @Param        page       query  int     false  "Page number"  default(1)
// @Param        page_size  query  int     false  "Page size"    default(10)
// @Success      200  {object}  models.PaginatedResponse
// @Failure      401  {object}  models.ErrorResponse
// @Router       /courses/{course_id}/lessons [get]
func (h *LessonHandler) ListByCourse(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("course_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid course ID")
		return
	}

	params := parsePagination(c)

	lessons, total, err := h.lessonService.ListByCourse(courseID, params)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:     lessons,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	})
}

// UpdateLesson godoc
// @Summary      Update a lesson (admin)
// @Description  Update lesson by ID
// @Tags         Lessons
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  int                          true  "Lesson ID"
// @Param        body  body  models.UpdateLessonRequest   true  "Update data"
// @Success      200   {object}  models.Lesson
// @Failure      400   {object}  models.ErrorResponse
// @Failure      401   {object}  models.ErrorResponse
// @Failure      403   {object}  models.ErrorResponse
// @Failure      404   {object}  models.ErrorResponse
// @Router       /lessons/{id} [put]
func (h *LessonHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("lesson_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid lesson ID")
		return
	}

	var req models.UpdateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err.Error())
		return
	}

	lesson, err := h.lessonService.Update(id, req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, lesson)
}

// DeleteLesson godoc
// @Summary      Delete a lesson (admin)
// @Description  Delete lesson by ID
// @Tags         Lessons
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Lesson ID"
// @Success      204
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /lessons/{id} [delete]
func (h *LessonHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("lesson_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid lesson ID")
		return
	}

	if err := h.lessonService.Delete(id); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
