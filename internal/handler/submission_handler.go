package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kaplenko/diplom/internal/models"
	"github.com/kaplenko/diplom/internal/service"
)

type SubmissionHandler struct {
	submissionService *service.SubmissionService
}

func NewSubmissionHandler(submissionService *service.SubmissionService) *SubmissionHandler {
	return &SubmissionHandler{submissionService: submissionService}
}

// CreateSubmission godoc
// @Summary      Submit a solution
// @Description  Submit code for a task
// @Tags         Submissions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        task_id  path  int                             true  "Task ID"
// @Param        body     body  models.CreateSubmissionRequest  true  "Submission data"
// @Success      201  {object}  models.Submission
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /tasks/{task_id}/submissions [post]
func (h *SubmissionHandler) Create(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("task_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid task ID")
		return
	}

	var req models.CreateSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err.Error())
		return
	}

	sub, err := h.submissionService.Create(taskID, getUserID(c), req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, sub)
}

// GetSubmission godoc
// @Summary      Get a submission
// @Description  Get submission by ID
// @Tags         Submissions
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Submission ID"
// @Success      200  {object}  models.Submission
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /submissions/{id} [get]
func (h *SubmissionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("submission_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid submission ID")
		return
	}

	sub, err := h.submissionService.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, sub)
}

// ListTaskSubmissions godoc
// @Summary      List own submissions for a task
// @Description  Get a paginated list of the authenticated user's submissions for a task
// @Tags         Submissions
// @Produce      json
// @Security     BearerAuth
// @Param        task_id    path   int  true   "Task ID"
// @Param        page       query  int  false  "Page number"  default(1)
// @Param        page_size  query  int  false  "Page size"    default(10)
// @Success      200  {object}  models.PaginatedResponse
// @Failure      401  {object}  models.ErrorResponse
// @Router       /tasks/{task_id}/submissions [get]
func (h *SubmissionHandler) ListByTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("task_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid task ID")
		return
	}

	params := parsePagination(c)

	subs, total, err := h.submissionService.ListByTask(taskID, getUserID(c), params)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:     subs,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	})
}

// ListAllSubmissions godoc
// @Summary      List all submissions (admin)
// @Description  Get a paginated list of all submissions
// @Tags         Submissions
// @Produce      json
// @Security     BearerAuth
// @Param        page       query  int  false  "Page number"  default(1)
// @Param        page_size  query  int  false  "Page size"    default(10)
// @Success      200  {object}  models.PaginatedResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Router       /submissions [get]
func (h *SubmissionHandler) ListAll(c *gin.Context) {
	params := parsePagination(c)

	subs, total, err := h.submissionService.ListAll(params)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:     subs,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	})
}
