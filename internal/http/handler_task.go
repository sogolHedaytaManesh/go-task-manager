package http

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"task-manager/internal/entities"
	"task-manager/internal/repository/postgres"
	"task-manager/pkg/rest"
	"time"
)

// TaskCreate handles the creation of a new task.
//
// @Summary Create a new task
// @Description Creates a new task with a title, description, status, and assignee.
//
//	If status is not provided, it defaults to "pending".
//
// @Tags Tasks
// @Accept json
// @Produce json
// @Param request body CreateTaskRequest true "Task creation payload"
// @Success 201 {object} rest.StandardResponse{data=TaskResponse} "Task successfully created"
// @Failure 400 {object} rest.StandardResponse{data=nil} "Invalid request payload or invalid task status"
// @Failure 500 {object} rest.StandardResponse{data=nil} "Internal server error"
// @Router /api/tasks/ [post]
func (h *Handler) TaskCreate(c *gin.Context) {
	traceID := TraceIDFromContext(c.Request.Context())

	h.logger.InfoF(LogTemplateIncoming, traceID, LogIncomingTaskCreate)

	var req CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorF(LogTemplateError, traceID, LogTaskCreateFailed, err.Error())
		c.JSON(http.StatusBadRequest, rest.GetFailedValidationResponse(err))
		return
	}

	if req.Status == "" {
		req.Status = entities.TaskStatusPending
	}

	if !req.Status.IsValid() {
		h.logger.ErrorF(LogTemplateError, traceID, LogTaskCreateFailed, errors.New(InvalidTaskStatus))
		c.JSON(http.StatusBadRequest, rest.GetFailedValidationResponse(errors.New(InvalidTaskStatus)))
		return
	}

	task := &entities.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		AssigneeID:  req.AssigneeID,
	}

	createdTask, err := h.TaskService.Create(c, task)
	if err != nil {
		h.logger.ErrorWithContext(c, fmt.Sprintf(LogTemplateError, traceID, Error, err.Error()))
		c.JSON(http.StatusInternalServerError, rest.InternalServerError)
		return
	}

	h.logger.InfoF(LogTemplateSuccess, traceID, LogTaskCreateSuccess, createdTask.ID)

	c.JSON(http.StatusCreated, rest.GetSuccessResponse(TaskResponse{
		ID:          createdTask.ID,
		Title:       createdTask.Title,
		Description: createdTask.Description,
		Status:      createdTask.Status,
		AssigneeID:  createdTask.AssigneeID,
		CreatedAt:   createdTask.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   createdTask.UpdatedAt.Format(time.RFC3339),
	}))
}

// TaskUpdate updates an existing task by its ID.
//
// @Summary Update an existing task
// @Description Updates the details of an existing task such as title, description, status, and assignee.
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body UpdateTaskRequest true "Updated task data"
// @Success 200 {object} rest.StandardResponse{data=TaskResponse} "Task successfully updated"
// @Failure 400 {object} rest.StandardResponse{data=nil} "Invalid task ID or request payload"
// @Failure 404 {object} rest.StandardResponse{data=nil} "Task not found"
// @Failure 500 {object} rest.StandardResponse{data=nil} "Internal server error"
// @Router /api/tasks/{id} [put]
func (h *Handler) TaskUpdate(c *gin.Context) {
	traceID := TraceIDFromContext(c.Request.Context())

	h.logger.InfoF(LogTemplateIncoming, traceID, LogIncomingTaskUpdate)
	taskIDParam := c.Param(ID)
	if taskIDParam == "" {
		h.logger.ErrorF(LogTemplateError, traceID, LogTaskUpdateFailed, errors.New(TaskIDIsRequired))
		c.JSON(http.StatusBadRequest, rest.GetFailedValidationResponse(errors.New(TaskIDIsRequired)))
		return
	}

	taskID, err := strconv.ParseInt(taskIDParam, 10, 64)
	if err != nil {
		h.logger.ErrorF(LogTemplateError, traceID, LogTaskUpdateFailed, errors.New(InvalidTaskID))
		c.JSON(http.StatusBadRequest, rest.GetFailedValidationResponse(errors.New(InvalidTaskStatus)))
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorF(LogTemplateError, traceID, LogTaskUpdateFailed, err.Error())
		c.JSON(http.StatusBadRequest, rest.GetFailedValidationResponse(err))
		return
	}

	task := &entities.Task{
		ID:          taskID,
		Title:       req.Title,
		Description: req.Description,
		AssigneeID:  req.AssigneeID,
		Status:      req.Status,
	}

	updatedTask, err := h.TaskService.Update(c, task)
	if err != nil {
		h.logger.ErrorWithContext(c, fmt.Sprintf(LogTemplateError, traceID, Error, err.Error()))
		if errors.Is(err, postgres.ErrTaskNotFound) {
			h.logger.ErrorF(LogTemplateError, traceID, LogTaskUpdateFailed, rest.NotFound)
			c.JSON(http.StatusNotFound, rest.NotFound)

			return
		}

		h.logger.ErrorF(LogTemplateError, traceID, LogTaskUpdateFailed, rest.InternalServerError)
		c.JSON(http.StatusInternalServerError, rest.InternalServerError)

		return
	}

	h.logger.InfoF(LogTemplateSuccess, traceID, LogTaskUpdateSuccess, updatedTask.ID)

	c.JSON(http.StatusOK, rest.GetSuccessResponse(TaskResponse{
		ID:          updatedTask.ID,
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		Status:      updatedTask.Status,
		AssigneeID:  updatedTask.AssigneeID,
		CreatedAt:   updatedTask.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   updatedTask.UpdatedAt.Format(time.RFC3339),
	}))
}

// TaskDelete deletes an existing task by its ID.
//
// @Summary Delete a task
// @Description Deletes an existing task identified by its ID.
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 204 "Task successfully deleted"
// @Failure 400 {object} rest.StandardResponse{data=nil} "Invalid task ID"
// @Failure 404 {object} rest.StandardResponse{data=nil} "Task not found"
// @Failure 500 {object} rest.StandardResponse{data=nil} "Internal server error"
// @Router /api/tasks/{id} [delete]
func (h *Handler) TaskDelete(c *gin.Context) {
	traceID := TraceIDFromContext(c.Request.Context())

	h.logger.InfoF(LogTemplateIncoming, traceID, LogIncomingTaskDelete)

	taskIDParam := c.Param(ID)
	if taskIDParam == "" {
		h.logger.ErrorF(LogTemplateError, traceID, LogTaskDeleteFailed, errors.New(TaskIDIsRequired))
		c.JSON(http.StatusBadRequest, rest.GetFailedValidationResponse(errors.New(TaskIDIsRequired)))
		return
	}

	taskID, err := strconv.ParseInt(taskIDParam, 10, 64)
	if err != nil {
		h.logger.ErrorF(LogTemplateError, traceID, LogTaskDeleteFailed, errors.New(InvalidTaskStatus))
		c.JSON(http.StatusBadRequest, rest.GetFailedValidationResponse(errors.New(InvalidTaskID)))
		return
	}

	err = h.TaskService.Delete(c, taskID)
	if err != nil {
		if errors.Is(err, postgres.ErrTaskNotFound) {
			h.logger.ErrorF(LogTemplateError, traceID, LogTaskDeleteFailed, rest.NotFound)
			c.JSON(http.StatusNotFound, rest.NotFound)
			return
		}

		h.logger.ErrorWithContext(c, fmt.Sprintf(LogTemplateError, traceID, Error, err.Error()))
		c.JSON(http.StatusInternalServerError, rest.InternalServerError)
		return
	}

	h.logger.InfoF(LogTemplateSuccess, traceID, LogTaskDeleteSuccess, taskID)

	c.Status(http.StatusNoContent)
}

// TaskGetByID retrieves a task by its ID.
//
// @Summary Get task by ID
// @Description Retrieves the details of a specific task using its ID.
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} rest.StandardResponse{data=TaskResponse} "Task successfully fetched"
// @Failure 400 {object} rest.StandardResponse{data=nil} "Invalid task ID"
// @Failure 404 {object} rest.StandardResponse{data=nil} "Task not found"
// @Failure 500 {object} rest.StandardResponse{data=nil} "Internal server error"
// @Router /api/tasks/{id} [get]
func (h *Handler) TaskGetByID(c *gin.Context) {
	traceID := TraceIDFromContext(c.Request.Context())

	h.logger.InfoF(LogTemplateIncoming, traceID, LogIncomingTaskFetch)

	taskIDParam := c.Param(ID)
	if taskIDParam == "" {
		h.logger.ErrorF(LogTemplateError, traceID, LogTaskFetchFailed, errors.New(TaskIDIsRequired))

		c.JSON(http.StatusBadRequest, rest.GetFailedValidationResponse(errors.New(TaskIDIsRequired)))

		return
	}

	taskID, err := strconv.ParseInt(taskIDParam, 10, 64)
	if err != nil {
		h.logger.ErrorF(LogTemplateError, traceID, LogTaskFetchFailed, errors.New(InvalidTaskID))
		c.JSON(http.StatusBadRequest, rest.GetFailedValidationResponse(errors.New(InvalidTaskID)))
		return
	}

	task, err := h.TaskService.GetByID(c, taskID)
	if err != nil {
		if errors.Is(err, postgres.ErrTaskNotFound) {
			h.logger.ErrorF(LogTemplateError, traceID, LogTaskFetchFailed, rest.NotFound)
			c.JSON(http.StatusNotFound, rest.NotFound)
			return
		}

		h.logger.ErrorWithContext(c, fmt.Sprintf(LogTemplateError, traceID, Error, err.Error()))
		c.JSON(http.StatusInternalServerError, rest.InternalServerError)
		return
	}

	h.logger.InfoF(LogTemplateSuccess, traceID, LogTaskFetchSuccess, taskID)

	c.JSON(http.StatusOK, rest.GetSuccessResponse(TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		AssigneeID:  task.AssigneeID,
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
	}))
}

// TaskList retrieves a paginated list of tasks.
//
// @Summary List tasks
// @Description Retrieves a paginated list of tasks with optional filters.
// @Tags Tasks
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Number of tasks per page" default(20)
// @Param filter query string false "Filter by any task field (e.g., title, status)"
// @Success 200 {object} rest.StandardResponse{data=[]entities.Task, meta=rest.PaginationMeta} "List of tasks successfully fetched"
// @Failure 500 {object} rest.StandardResponse{data=nil} "Internal server error"
// @Router /api/tasks [get]
func (h *Handler) TaskList(c *gin.Context) {
	traceID := TraceIDFromContext(c.Request.Context())

	h.logger.InfoF(LogTemplateIncoming, traceID, LogIncomingTaskFetch)

	query := rest.ParseQuery(c)

	tasks, total, err := h.TaskService.List(c, query)
	if err != nil {
		h.logger.ErrorWithContext(c, fmt.Sprintf(LogTemplateError, traceID, Error, err.Error()))

		c.JSON(http.StatusInternalServerError, rest.InternalServerError)

		return
	}

	h.logger.InfoF(LogTemplateSuccess, traceID, LogTaskFetchSuccess, Bulk)

	c.JSON(http.StatusOK, rest.GetSuccessResponseWithMeta(tasks, rest.PaginationMeta{
		Total:   total,
		Page:    query.Page,
		PerPage: query.PerPage,
	},
	),
	)
}

const (
	LogIncomingTaskCreate = "Incoming task create request"
	LogTaskCreateSuccess  = "Task created successfully"
	LogTaskCreateFailed   = "Failed to create task"

	LogIncomingTaskUpdate = "Incoming task update request"
	LogTaskUpdateSuccess  = "Task updated successfully"
	LogTaskUpdateFailed   = "Failed to update task"

	LogIncomingTaskDelete = "Incoming task delete request"
	LogTaskDeleteSuccess  = "Task delete successfully"
	LogTaskDeleteFailed   = "Failed to delete task"

	LogIncomingTaskFetch = "Incoming task fetch request"
	LogTaskFetchSuccess  = "Task fetch successfully"
	LogTaskFetchFailed   = "Failed to fetch task"

	InvalidTaskStatus = "Invalid task status"
	TaskIDIsRequired  = "Task ID is required"
	InvalidTaskID     = "Invalid task ID"

	ID    = "id"
	Error = "err"
	Bulk  = "bulk"
)

var (
	LogTemplateIncoming = "[TRACE %s] %s"
	LogTemplateSuccess  = "[TRACE %s] %s: ID=%d"
	LogTemplateError    = "[TRACE %s] %s: %v"
)

type CreateTaskRequest struct {
	Title       string              `json:"title" binding:"required"`
	Description string              `json:"description,omitempty"`
	Status      entities.TaskStatus `json:"status,omitempty"` // optional, default pending
	AssigneeID  int64               `json:"assignee_id" binding:"required"`
}

type UpdateTaskRequest struct {
	Title       string              `json:"title" binding:"required"`
	Description string              `json:"description,omitempty"`
	Status      entities.TaskStatus `json:"status" binding:"required,oneof=pending in_progress done canceled"`
	AssigneeID  int64               `json:"assignee_id" binding:"required"`
}

type TaskResponse struct {
	ID          int64               `json:"id"`
	Title       string              `json:"title"`
	Description string              `json:"description,omitempty"`
	Status      entities.TaskStatus `json:"status"`
	AssigneeID  int64               `json:"assignee_id" binding:"required"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
}
