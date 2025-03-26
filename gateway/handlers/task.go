package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"sync"
	"log"
	"context"

	"sama/go-task-management/commons"
	"sama/go-task-management/gateway/handlers/constants"
	"sama/go-task-management/gateway/handlers/validation"
	"sama/go-task-management/gateway/middleware"
	"sama/go-task-management/gateway/services/auth"
	"sama/go-task-management/gateway/services/task"

	"github.com/google/uuid"
)

type TaskHandler struct {
	*BaseHandler
	taskService *task.Service
}

func NewTaskHandler(base *BaseHandler, taskService *task.Service) *TaskHandler {
	return &TaskHandler{
		BaseHandler: base,
		taskService: taskService,
	}
}

type GetTaskResponse struct {
	ID          string                    `json:"id"`
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	Status      string                    `json:"status"`
	Priority    int                       `json:"priority"`
	DueDate     time.Time                 `json:"due_date"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
	Creator     auth.UserResponse         `json:"creator"`
	Assignee    *auth.UserResponse        `json:"assignee,omitempty"`
	Events      []TaskSystemEventResponse `json:"events"`
}

type GetAllTasksResponse struct {
	Tasks []GetTaskResponse `json:"tasks"`
}

type CreateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    int        `json:"priority"`
	DueDate     time.Time  `json:"due_date"`
	AssigneeID  *string    `json:"assignee_id,omitempty"`
}

func (r *CreateTaskRequest) Validate() []validation.ValidationError {
	var errors []validation.ValidationError

	if titleErr := validation.ValidateTitle(r.Title); titleErr != nil {
		errors = append(errors, *titleErr)
	}

	if descErr := validation.ValidateDescription(r.Description); descErr != nil {
		errors = append(errors, *descErr)
	}

	if r.Status != constants.TaskStatusTodo && r.Status != constants.TaskStatusInProgress && r.Status != constants.TaskStatusDone {
		errors = append(errors, validation.ValidationError{
			Field:   "status",
			Message: "Invalid status. Must be one of: TODO, IN_PROGRESS, DONE",
		})
	}

	if r.Priority < constants.TaskPriorityLow || r.Priority > constants.TaskPriorityHigh {
		errors = append(errors, validation.ValidationError{
			Field:   "priority",
			Message: "Priority must be between 1 and 3",
		})
	}

	if r.DueDate.Before(time.Now()) {
		errors = append(errors, validation.ValidationError{
			Field:   "due_date",
			Message: "Due date must be in the future",
		})
	}

	return errors
}

type CreateTaskResponse struct {
	TaskId string `json:"task_id"`
}

type UpdateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    int        `json:"priority"`
	DueDate     time.Time  `json:"due_date"`
	AssigneeID  *string    `json:"assignee_id,omitempty"`
}

func (r *UpdateTaskRequest) Validate() []validation.ValidationError {
	var errors []validation.ValidationError

	if r.Title != "" {
		if titleErr := validation.ValidateTitle(r.Title); titleErr != nil {
			errors = append(errors, *titleErr)
		}
	}

	if r.Description != "" {
		if descErr := validation.ValidateDescription(r.Description); descErr != nil {
			errors = append(errors, *descErr)
		}
	}

	if r.Status != "" && r.Status != constants.TaskStatusTodo && r.Status != constants.TaskStatusInProgress && r.Status != constants.TaskStatusDone {
		errors = append(errors, validation.ValidationError{
			Field:   "status",
			Message: "Invalid status. Must be one of: TODO, IN_PROGRESS, DONE",
		})
	}

	if r.Priority != 0 && (r.Priority < constants.TaskPriorityLow || r.Priority > constants.TaskPriorityHigh) {
		errors = append(errors, validation.ValidationError{
			Field:   "priority",
			Message: "Priority must be between 1 and 3",
		})
	}

	if !r.DueDate.IsZero() && r.DueDate.Before(time.Now()) {
		errors = append(errors, validation.ValidationError{
			Field:   "due_date",
			Message: "Due date must be in the future",
		})
	}

	return errors
}

type UpdateTaskResponse struct {
	TaskId string `json:"task_id"`
}

// @Summary Get a task by ID
// @Description Retrieves a specific task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} GetTaskResponse
// @Failure 400 {object} ErrorResponse "Invalid task ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Task not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")
	if taskID == "" {
		h.respondWithError(w, http.StatusBadRequest, constants.ErrCodeBadRequest, "Task ID is required", "")
		return
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Unauthorized", "")
		return
	}

	task, err := h.taskService.GetTask(r.Context(), taskID, userID)
	if err != nil {
		switch err {
		case commons.ErrNotFound:
			h.respondWithError(w, http.StatusNotFound, constants.ErrCodeNotFound, "Task not found", "")
		case commons.ErrUnauthorized:
			h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Unauthorized", "")
		default:
			h.respondWithError(w, http.StatusInternalServerError, constants.ErrCodeInternal, "Internal server error", "")
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, StandardResponse{
		Success: true,
		Data:    task,
	})
}

// @Summary Get all tasks
// @Description Retrieves all tasks for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} GetAllTasksResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /tasks [get]
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Unauthorized", "")
		return
	}

	tasks, err := h.taskService.GetAllTasks(r.Context(), userID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, constants.ErrCodeInternal, "Failed to fetch tasks", err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, StandardResponse{
		Success: true,
		Data:    tasks,
	})
}

// @Summary Create a new task
// @Description Creates a new task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param input body CreateTaskRequest true "Task details"
// @Success 201 {object} CreateTaskResponse
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input task.CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondWithError(w, http.StatusBadRequest, constants.ErrCodeBadRequest, "Invalid request payload", err.Error())
		return
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Unauthorized", "")
		return
	}

	input.CreatorID = userID

	task, err := h.taskService.CreateTask(r.Context(), input)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, constants.ErrCodeInternal, "Failed to create task", err.Error())
		return
	}

	correlationId := uuid.New().String()
	var wg sync.WaitGroup

  wg.Add(1)
  go func() {
    defer wg.Done()

		_, err := h.taskEventService.Create(
			task.ID,
			correlationId,
			"API Gateway",
      "api:request:received",
      "Task creation request received",
			"{}", //marshallJson(params),
			0,
		)

    if err != nil {
      log.Printf("Error creating request event: %v", err)
    }
  }()

	wg.Add(1)
  go func() {
    defer wg.Done()

		_, err := h.taskEventService.Create(
			task.ID,
			correlationId,
			"API Gateway",
      "api:db:task-created",
      "Task created in database",
			"{}",
			1,
		)

    if err != nil {
      log.Printf("Error creating task created event: %v", err)
    }
  }()

	wg.Add(1)
  go func() {
    defer wg.Done()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

		grpcErr := h.grpcService.SendNotification(ctx, commons.GRPCEvent{
			TaskId:        task.ID,
			CorrelationId: correlationId,
			Types:         []string{"IN_APP", "EMAIL"},
		})
		if grpcErr != nil {
			log.Printf("Failed to send notification: %v", grpcErr)
		}

		_, err := h.taskEventService.Create(
			task.ID,
			correlationId,
			"API Gateway",
      "api:event:task-created",
      "Task created event emitted",
			"{}",
			2,
		)

		if err != nil {
      log.Printf("Failed to create task created event: %v", err)
    }
  }()

	wg.Wait()

	h.respondWithJSON(w, http.StatusCreated, StandardResponse{
		Success: true,
		Data:    task,
	})
}

// @Summary Update a task
// @Description Updates an existing task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param input body UpdateTaskRequest true "Task update details"
// @Success 200 {object} UpdateTaskResponse
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Task not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")
	if taskID == "" {
		h.respondWithError(w, http.StatusBadRequest, constants.ErrCodeBadRequest, "Task ID is required", "")
		return
	}

	var input task.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondWithError(w, http.StatusBadRequest, constants.ErrCodeBadRequest, "Invalid request payload", err.Error())
		return
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Unauthorized", "")
		return
	}

	input.UserID = userID

	task, err := h.taskService.UpdateTask(r.Context(), taskID, input)
	if err != nil {
		switch err {
		case commons.ErrNotFound:
			h.respondWithError(w, http.StatusNotFound, constants.ErrCodeNotFound, "Task not found", "")
		case commons.ErrUnauthorized:
			h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Unauthorized", "")
		default:
			h.respondWithError(w, http.StatusInternalServerError, constants.ErrCodeInternal, "Failed to update task", err.Error())
		}
		return
	}

	correlationId := uuid.New().String()
	_, errEvent := h.taskEventService.Create(
		taskID,
		correlationId,
		"API Gateway",
    "api:event:task-updated",
    "Task updated event emitted",
    "{}",
		3,
	)
	if errEvent != nil {
		log.Printf("Failed to create task updated event: %v", errEvent)
	}

	h.respondWithJSON(w, http.StatusOK, StandardResponse{
		Success: true,
		Data:    task,
	})
}

// @Summary Delete a task
// @Description Deletes a task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 "Task deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid task ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Task not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")
	if taskID == "" {
		h.respondWithError(w, http.StatusBadRequest, constants.ErrCodeBadRequest, "Task ID is required", "")
		return
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Unauthorized", "")
		return
	}

	err := h.taskService.DeleteTask(r.Context(), taskID, userID)
	if err != nil {
		switch err {
		case commons.ErrNotFound:
			h.respondWithError(w, http.StatusNotFound, constants.ErrCodeNotFound, "Task not found", "")
		case commons.ErrUnauthorized:
			h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Unauthorized", "")
		default:
			h.respondWithError(w, http.StatusInternalServerError, constants.ErrCodeInternal, "Failed to delete task", err.Error())
		}
		return
	}

	correlationId := uuid.New().String()
	_, errEvent := h.taskEventService.Create(
		taskID,
		correlationId,
		"API Gateway",
    "api:event:task-deleted",
    "Task deleted event emitted",
    "{}",
		4,
	)
	if errEvent != nil {
		log.Printf("Failed to create task deleted event: %v", errEvent)
	}

	h.respondWithJSON(w, http.StatusOK, StandardResponse{
		Success: true,
		Data: map[string]string{
			"message": "Task deleted successfully",
		},
	})
}
