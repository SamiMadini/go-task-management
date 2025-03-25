package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"sama/go-task-management/commons"
	pb "sama/go-task-management/commons/api"
	"sama/go-task-management/gateway/middleware"

	"github.com/google/uuid"
)

const notificationTimeout = 5 * time.Second

type TaskHandler struct {
	*BaseHandler
}

func NewTaskHandler(base *BaseHandler) *TaskHandler {
	return &TaskHandler{BaseHandler: base}
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
	Creator     UserResponse              `json:"creator"`
	Assignee    *UserResponse             `json:"assignee,omitempty"`
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

func (r *CreateTaskRequest) Validate() []ValidationError {
	var errors []ValidationError

	if r.Title == "" {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "Title is required",
		})
	}

	if len(r.Title) > 200 {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "Title must be less than 200 characters",
		})
	}

	if len(r.Description) > 1000 {
		errors = append(errors, ValidationError{
			Field:   "description",
			Message: "Description must be less than 1000 characters",
		})
	}

	if r.Status != string(TaskStatusTodo) && r.Status != string(TaskStatusInProgress) && r.Status != string(TaskStatusDone) {
		errors = append(errors, ValidationError{
			Field:   "status",
			Message: "Invalid status. Must be one of: TODO, IN_PROGRESS, DONE",
		})
	}

	if r.Priority < int(TaskPriorityLow) || r.Priority > int(TaskPriorityHigh) {
		errors = append(errors, ValidationError{
			Field:   "priority",
			Message: "Priority must be between 1 and 3",
		})
	}

	if r.DueDate.Before(time.Now()) {
		errors = append(errors, ValidationError{
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

func (r *UpdateTaskRequest) Validate() []ValidationError {
	var errors []ValidationError

	if r.Title != "" && len(r.Title) > 200 {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "Title must be less than 200 characters",
		})
	}

	if r.Description != "" && len(r.Description) > 1000 {
		errors = append(errors, ValidationError{
			Field:   "description",
			Message: "Description must be less than 1000 characters",
		})
	}

	if r.Status != "" && r.Status != string(TaskStatusTodo) && r.Status != string(TaskStatusInProgress) && r.Status != string(TaskStatusDone) {
		errors = append(errors, ValidationError{
			Field:   "status",
			Message: "Invalid status. Must be one of: TODO, IN_PROGRESS, DONE",
		})
	}

	if r.Priority != 0 && (r.Priority < int(TaskPriorityLow) || r.Priority > int(TaskPriorityHigh)) {
		errors = append(errors, ValidationError{
			Field:   "priority",
			Message: "Priority must be between 1 and 3",
		})
	}

	if !r.DueDate.IsZero() && r.DueDate.Before(time.Now()) {
		errors = append(errors, ValidationError{
			Field:   "due_date",
			Message: "Due date must be in the future",
		})
	}

	return errors
}

type UpdateTaskResponse struct {
	TaskId string `json:"task_id"`
}

func createTaskSystemEvent(taskId, correlationId, origin, action, message, jsonData string) commons.TaskSystemEvent {
	now := time.Now()
	return commons.TaskSystemEvent{
		ID:            uuid.New().String(),
		TaskId:        taskId,
		CorrelationId: correlationId,
		Origin:        origin,
		Action:        action,
		Message:       message,
		JsonData:      jsonData,
		EmitAt:        now,
		CreatedAt:     now,
	}
}

func (h *TaskHandler) marshallJson(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		h.logger.Printf("Error marshalling JSON: %v", err)
		return "{}"
	}
	return string(data)
}

// @Summary Get a task by ID
// @Description Retrieves a specific task by its ID
// @Description
// @Description Error scenarios:
// @Description - Invalid UUID format: Returns 400 with BAD_REQUEST code
// @Description - Task not found: Returns 404 with NOT_FOUND code
// @Description - Unauthorized access: Returns 403 with FORBIDDEN code
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID (UUID format)"
// @Success 200 {object} GetTaskResponse
// @Failure 400 {object} ErrorResponse "Invalid task ID format"
// @Failure 404 {object} ErrorResponse "Task not found"
// @Failure 403 {object} ErrorResponse "Access denied - not task owner or assignee"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Task ID is required", "Path parameter 'id' is missing")
		return
	}

	if !commons.IsValidUUID(id) {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid task ID format", "Task ID must be a valid UUID")
		return
	}

	task, err := h.taskRepository.GetByID(id)
	if err != nil {
		h.logger.Printf("Error fetching task %s: %v", id, err)
		h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to fetch task", err.Error())
		return
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID != task.CreatorID && (task.AssigneeID == nil || *task.AssigneeID != userID) {
		h.respondWithError(w, http.StatusForbidden, ErrCodeForbidden, "Access denied", "You don't have permission to view this task")
		return
	}

	creator, err := h.userRepository.GetByID(task.CreatorID)
	if err != nil {
		h.logger.Printf("Error fetching creator details: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to fetch creator details", err.Error())
		return
	}

	var assignee *commons.User
	if task.AssigneeID != nil {
		assigneeUser, err := h.userRepository.GetByID(*task.AssigneeID)
		if err != nil {
			h.logger.Printf("Error fetching assignee details: %v", err)
			h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to fetch assignee details", err.Error())
			return
		}
		assignee = &assigneeUser
	}

	events := make([]TaskSystemEventResponse, len(task.Events))
	for i, event := range task.Events {
		events[i] = TaskSystemEventResponse{
			ID:            event.ID,
			TaskId:        event.TaskId,
			CorrelationId: event.CorrelationId,
			Origin:        event.Origin,
			Action:        event.Action,
			Message:       event.Message,
			JsonData:      event.JsonData,
			EmitAt:        event.EmitAt,
			CreatedAt:     event.CreatedAt,
		}
	}

	response := GetTaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Creator:     UserResponse{ID: creator.ID, Handle: creator.Handle, Email: creator.Email, Status: creator.Status},
		Events:      events,
	}

	if assignee != nil {
		response.Assignee = &UserResponse{ID: assignee.ID, Handle: assignee.Handle, Email: assignee.Email, Status: assignee.Status}
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// @Summary Get all tasks
// @Description Retrieves all tasks for the current user (created or assigned)
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} GetAllTasksResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks [get]
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	tasks, err := h.taskRepository.GetByUserID(userID)
	if err != nil {
		h.logger.Printf("Failed to get tasks for user %s: %v", userID, err)
		h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to fetch tasks", err.Error())
		return
	}

	if len(tasks) == 0 {
		h.respondWithJSON(w, http.StatusOK, GetAllTasksResponse{
			Tasks: []GetTaskResponse{},
		})
		return
	}

	response := GetAllTasksResponse{
		Tasks: make([]GetTaskResponse, len(tasks)),
	}

	for i, t := range tasks {
		// Get creator details
		creator, err := h.userRepository.GetByID(t.CreatorID)
		if err != nil {
			h.logger.Printf("Failed to get creator details for task %s: %v", t.ID, err)
			h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to fetch creator details", err.Error())
			return
		}

		// Get assignee details if exists
		var assignee *commons.User
		if t.AssigneeID != nil {
			assigneeUser, err := h.userRepository.GetByID(*t.AssigneeID)
			if err != nil {
				h.logger.Printf("Failed to get assignee details for task %s: %v", t.ID, err)
				h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to fetch assignee details", err.Error())
				return
			}
			assignee = &assigneeUser
		}

		events := make([]TaskSystemEventResponse, len(t.Events))
		for j, event := range t.Events {
			events[j] = TaskSystemEventResponse{
				ID:            event.ID,
				TaskId:        event.TaskId,
				CorrelationId: event.CorrelationId,
				Origin:        event.Origin,
				Action:        event.Action,
				Message:       event.Message,
				JsonData:      event.JsonData,
				EmitAt:        event.EmitAt,
				CreatedAt:     event.CreatedAt,
			}
		}

		response.Tasks[i] = GetTaskResponse{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Status:      t.Status,
			Priority:    t.Priority,
			DueDate:     t.DueDate,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
			Creator: UserResponse{
				ID:     creator.ID,
				Handle: creator.Handle,
				Email:  creator.Email,
				Status: creator.Status,
			},
			Events: events,
		}

		if assignee != nil {
			response.Tasks[i].Assignee = &UserResponse{
				ID:     assignee.ID,
				Handle: assignee.Handle,
				Email:  assignee.Email,
				Status: assignee.Status,
			}
		}
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// @Summary Create a new task
// @Description Creates a new task in the system
// @Description
// @Description Error scenarios:
// @Description - Missing required fields: Returns 400 with VALIDATION_ERROR code
// @Description - Invalid assignee ID: Returns 400 with BAD_REQUEST code
// @Description - Assignee not found: Returns 404 with NOT_FOUND code
// @Description - Invalid task status: Returns 400 with VALIDATION_ERROR code
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body CreateTaskRequest true "Task details"
// @Success 201 {object} CreateTaskResponse
// @Failure 400 {object} ErrorResponse "Validation errors or invalid input"
// @Failure 401 {object} ErrorResponse "Missing or invalid authentication token"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body", err.Error())
		return
	}

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		h.respondWithValidationErrors(w, validationErrors)
		return
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	if req.AssigneeID != nil {
		assignee, err := h.userRepository.GetByID(*req.AssigneeID)
		if err != nil {
			h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid assignee", "Assignee not found")
			return
		}
		if assignee.Status != "ACTIVE" {
			h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid assignee", "Assignee is not active")
			return
		}
	}

	task := commons.Task{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		CreatorID:   userID,
		AssigneeID:  req.AssigneeID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdTask, err := h.taskRepository.Create(task)
	if err != nil {
		h.logger.Printf("Error creating task: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to create task", err.Error())
		return
	}

	event := createTaskSystemEvent(
		createdTask.ID,
		uuid.New().String(),
		"API",
		"CREATE",
		"Task created",
		h.marshallJson(createdTask),
	)

	_, err = h.taskSystemEventRepository.Create(event, 1)
	if err != nil {
		h.logger.Printf("Error creating task system event: %v", err)
	}

	if createdTask.AssigneeID != nil {
		ctx, cancel := context.WithTimeout(context.Background(), notificationTimeout)
		defer cancel()

		notification := &pb.SendNotificationRequest{
			TaskId:        createdTask.ID,
			CorrelationId: uuid.New().String(),
			Types:         []pb.NotificationType{0},
		}

		if _, err := h.notificationServiceClient.SendNotification(ctx, notification); err != nil {
			h.logger.Printf("Error sending notification: %v", err)
		}
	}

	h.respondWithJSON(w, http.StatusCreated, CreateTaskResponse{TaskId: createdTask.ID})
}

// @Summary Update a task
// @Description Updates an existing task
// @Description
// @Description Error scenarios:
// @Description - Invalid UUID format: Returns 400 with BAD_REQUEST code
// @Description - Task not found: Returns 404 with NOT_FOUND code
// @Description - Invalid status transition: Returns 400 with VALIDATION_ERROR code
// @Description - Unauthorized modification: Returns 403 with FORBIDDEN code
// @Description - Invalid assignee: Returns 400 with BAD_REQUEST code
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID (UUID format)"
// @Param task body UpdateTaskRequest true "Updated task details"
// @Success 200 {object} UpdateTaskResponse
// @Failure 400 {object} ErrorResponse "Invalid input or validation errors"
// @Failure 401 {object} ErrorResponse "Missing or invalid authentication token"
// @Failure 403 {object} ErrorResponse "Not authorized to update this task"
// @Failure 404 {object} ErrorResponse "Task not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Task ID is required", "Path parameter 'id' is missing")
		return
	}

	if !commons.IsValidUUID(id) {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid task ID format", "Task ID must be a valid UUID")
		return
	}

	var req UpdateTaskRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body", err.Error())
		return
	}

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		h.respondWithValidationErrors(w, validationErrors)
		return
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	task, err := h.taskRepository.GetByID(id)
	if err != nil {
		h.logger.Printf("Error fetching task %s: %v", id, err)
		h.respondWithError(w, http.StatusNotFound, ErrCodeNotFound, "Task not found", err.Error())
		return
	}

	if userID != task.CreatorID && (task.AssigneeID == nil || *task.AssigneeID != userID) {
		h.respondWithError(w, http.StatusForbidden, ErrCodeForbidden, "Access denied", "You don't have permission to update this task")
		return
	}

	now := time.Now()
	task.Title = req.Title
	task.Description = req.Description
	task.Status = req.Status
	task.Priority = req.Priority
	task.DueDate = req.DueDate
	task.AssigneeID = req.AssigneeID
	task.UpdatedAt = now

	if err := h.taskRepository.Update(task); err != nil {
		h.logger.Printf("Error updating task %s: %v", id, err)
		h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to update task", err.Error())
		return
	}

	correlationId := uuid.New().String()
	event := createTaskSystemEvent(
		id,
		correlationId,
		"API",
		"UPDATE",
		"Task updated",
		h.marshallJson(task),
	)

	_, err = h.taskSystemEventRepository.Create(event, 1)
	if err != nil {
		h.logger.Printf("Error creating task updated event: %v", err)
	}

	h.respondWithJSON(w, http.StatusOK, UpdateTaskResponse{
		TaskId: id,
	})
}

// @Summary Delete a task
// @Description Deletes a task from the system
// @Description
// @Description Error scenarios:
// @Description - Invalid UUID format: Returns 400 with BAD_REQUEST code
// @Description - Task not found: Returns 404 with NOT_FOUND code
// @Description - Unauthorized deletion: Returns 403 with FORBIDDEN code (only creator can delete)
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID (UUID format)"
// @Success 200 {object} map[string]bool "Success status"
// @Failure 400 {object} ErrorResponse "Invalid task ID format"
// @Failure 401 {object} ErrorResponse "Missing or invalid authentication token"
// @Failure 403 {object} ErrorResponse "Not authorized to delete this task"
// @Failure 404 {object} ErrorResponse "Task not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Task ID is required", "Path parameter 'id' is missing")
		return
	}

	if !commons.IsValidUUID(id) {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid task ID format", "Task ID must be a valid UUID")
		return
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	task, err := h.taskRepository.GetByID(id)
	if err != nil {
		h.logger.Printf("Error fetching task %s: %v", id, err)
		h.respondWithError(w, http.StatusNotFound, ErrCodeNotFound, "Task not found", err.Error())
		return
	}

	if userID != task.CreatorID {
		h.respondWithError(w, http.StatusForbidden, ErrCodeForbidden, "Access denied", "Only the creator can delete a task")
		return
	}

	if err := h.taskRepository.Delete(id); err != nil {
		h.logger.Printf("Error deleting task %s: %v", id, err)
		h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to delete task", err.Error())
		return
	}

	correlationId := uuid.New().String()
	event := createTaskSystemEvent(
		id,
		correlationId,
		"API",
		"DELETE",
		"Task deleted",
		h.marshallJson(task),
	)

	_, err = h.taskSystemEventRepository.Create(event, 1)
	if err != nil {
		h.logger.Printf("Error creating task deleted event: %v", err)
	}

	h.respondWithJSON(w, http.StatusOK, map[string]bool{
		"success": true,
	})
}
