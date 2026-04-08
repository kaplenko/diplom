package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kaplenko/diplom/internal/models"
	"github.com/kaplenko/diplom/internal/service"
)

type CourseHandler struct {
	courseService *service.CourseService
}

func NewCourseHandler(courseService *service.CourseService) *CourseHandler {
	return &CourseHandler{courseService: courseService}
}

// CreateCourse godoc
// @Summary      Create a course (admin)
// @Description  Create a new course
// @Tags         Courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      models.CreateCourseRequest  true  "Course data"
// @Success      201   {object}  models.Course
// @Failure      400   {object}  models.ErrorResponse
// @Failure      401   {object}  models.ErrorResponse
// @Failure      403   {object}  models.ErrorResponse
// @Router       /courses [post]
func (h *CourseHandler) Create(c *gin.Context) {
	var req models.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err.Error())
		return
	}

	course, err := h.courseService.Create(req, getUserID(c))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, course)
}

// GetCourse godoc
// @Summary      Get a course
// @Description  Get course by ID
// @Tags         Courses
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Course ID"
// @Success      200  {object}  models.Course
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /courses/{id} [get]
func (h *CourseHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("course_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid course ID")
		return
	}

	course, err := h.courseService.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, course)
}

// ListCourses godoc
// @Summary      List courses
// @Description  Get a paginated list of courses
// @Tags         Courses
// @Produce      json
// @Security     BearerAuth
// @Param        page       query  int     false  "Page number"  default(1)
// @Param        page_size  query  int     false  "Page size"    default(10)
// @Param        search     query  string  false  "Search by title or description"
// @Success      200  {object}  models.PaginatedResponse
// @Failure      401  {object}  models.ErrorResponse
// @Router       /courses [get]
func (h *CourseHandler) List(c *gin.Context) {
	params := parsePagination(c)

	courses, total, err := h.courseService.List(params)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:     courses,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	})
}

// UpdateCourse godoc
// @Summary      Update a course (admin)
// @Description  Update course by ID
// @Tags         Courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  int                         true  "Course ID"
// @Param        body  body  models.UpdateCourseRequest  true  "Update data"
// @Success      200   {object}  models.Course
// @Failure      400   {object}  models.ErrorResponse
// @Failure      401   {object}  models.ErrorResponse
// @Failure      403   {object}  models.ErrorResponse
// @Failure      404   {object}  models.ErrorResponse
// @Router       /courses/{id} [put]
func (h *CourseHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("course_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid course ID")
		return
	}

	var req models.UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err.Error())
		return
	}

	course, err := h.courseService.Update(id, req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, course)
}

// DeleteCourse godoc
// @Summary      Delete a course (admin)
// @Description  Delete course by ID
// @Tags         Courses
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Course ID"
// @Success      204
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /courses/{id} [delete]
func (h *CourseHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("course_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid course ID")
		return
	}

	if err := h.courseService.Delete(id); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
