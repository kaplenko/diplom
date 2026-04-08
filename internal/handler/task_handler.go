package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kaplenko/diplom/internal/models"
	"github.com/kaplenko/diplom/internal/service"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

// CreateTask godoc
// @Summary      Create a task (admin)
// @Description  Create a new task in a lesson
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        lesson_id  path  int                        true  "Lesson ID"
// @Param        body       body  models.CreateTaskRequest   true  "Task data"
// @Success      201  {object}  models.Task
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /lessons/{lesson_id}/tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	lessonID, err := strconv.ParseInt(c.Param("lesson_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid lesson ID")
		return
	}

	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err.Error())
		return
	}

	task, err := h.taskService.Create(lessonID, req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetTask godoc
// @Summary      Get a task
// @Description  Get task by ID
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Task ID"
// @Success      200  {object}  models.Task
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /tasks/{id} [get]
func (h *TaskHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("task_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid task ID")
		return
	}

	task, err := h.taskService.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

// ListTasks godoc
// @Summary      List tasks in a lesson
// @Description  Get a paginated list of tasks for a lesson
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        lesson_id  path   int     true   "Lesson ID"
// @Param        page       query  int     false  "Page number"  default(1)
// @Param        page_size  query  int     false  "Page size"    default(10)
// @Success      200  {object}  models.PaginatedResponse
// @Failure      401  {object}  models.ErrorResponse
// @Router       /lessons/{lesson_id}/tasks [get]
func (h *TaskHandler) ListByLesson(c *gin.Context) {
	lessonID, err := strconv.ParseInt(c.Param("lesson_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid lesson ID")
		return
	}

	params := parsePagination(c)

	tasks, total, err := h.taskService.ListByLesson(lessonID, params)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:     tasks,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	})
}

// UpdateTask godoc
// @Summary      Update a task (admin)
// @Description  Update task by ID
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  int                        true  "Task ID"
// @Param        body  body  models.UpdateTaskRequest   true  "Update data"
// @Success      200   {object}  models.Task
// @Failure      400   {object}  models.ErrorResponse
// @Failure      401   {object}  models.ErrorResponse
// @Failure      403   {object}  models.ErrorResponse
// @Failure      404   {object}  models.ErrorResponse
// @Router       /tasks/{id} [put]
func (h *TaskHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("task_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid task ID")
		return
	}

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err.Error())
		return
	}

	task, err := h.taskService.Update(id, req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask godoc
// @Summary      Delete a task (admin)
// @Description  Delete task by ID
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Task ID"
// @Success      204
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("task_id"), 10, 64)
	if err != nil {
		respondValidationError(c, "invalid task ID")
		return
	}

	if err := h.taskService.Delete(id); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
