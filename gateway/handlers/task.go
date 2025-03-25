package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"sama/go-task-management/commons"
	pb "sama/go-task-management/commons/api"
	"sama/go-task-management/gateway/middleware"

	"github.com/google/uuid"
	"google.golang.org/grpc"
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

func (r *CreateTaskRequest) Validate() error {
	// TODO: Add validation
	return nil
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

func (r *UpdateTaskRequest) Validate() error {
	// TODO: Add validation
	return nil
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

func marshallJson(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return "{}"
	}
	return string(data)
}

// @Summary Get a task by ID
// @Description Retrieves a specific task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} GetTaskResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.respondWithError(w, http.StatusBadRequest, "task ID is required")
		return
	}

	if !commons.IsValidUUID(id) {
		h.respondWithError(w, http.StatusBadRequest, "invalid task ID format")
		return
	}

	task, err := h.taskRepository.GetByID(id)
	if err != nil {
		log.Printf("Error fetching task %s: %v", id, err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch task")
		return
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID != task.CreatorID && (task.AssigneeID == nil || *task.AssigneeID != userID) {
		h.respondWithError(w, http.StatusForbidden, "You don't have permission to view this task")
		return
	}

	creator, err := h.userRepository.GetByID(task.CreatorID)
	if err != nil {
		log.Printf("Error fetching creator details: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch creator details")
		return
	}

	var assignee *commons.User
	if task.AssigneeID != nil {
		assigneeUser, err := h.userRepository.GetByID(*task.AssigneeID)
		if err != nil {
			log.Printf("Error fetching assignee details: %v", err)
			h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch assignee details")
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
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	tasks, err := h.taskRepository.GetByUserID(userID)
	if err != nil {
		log.Printf("Failed to get tasks for user %s: %v", userID, err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch tasks")
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
			log.Printf("Failed to get creator details for task %s: %v", t.ID, err)
			h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch creator details")
			return
		}

		// Get assignee details if exists
		var assignee *commons.User
		if t.AssigneeID != nil {
			assigneeUser, err := h.userRepository.GetByID(*t.AssigneeID)
			if err != nil {
				log.Printf("Failed to get assignee details for task %s: %v", t.ID, err)
				h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch assignee details")
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
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body CreateTaskRequest true "Task details"
// @Success 201 {object} CreateTaskResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var params CreateTaskRequest
	if err := h.decodeJSON(r, &params); err != nil {
		log.Printf("Invalid task creation request: %v", err)
		h.respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	now := time.Now()
	taskId := uuid.New().String()
	correlationId := uuid.New().String()

	task := commons.Task{
		ID:          taskId,
		CreatorID:   userID,
		AssigneeID:  params.AssigneeID,
		Title:       params.Title,
		Description: params.Description,
		Status:      params.Status,
		Priority:    params.Priority,
		DueDate:     params.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	var result commons.Task
	result, err := h.taskRepository.Create(task)
	if err != nil {
		log.Printf("Error creating task: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}
	task = result

	var wg sync.WaitGroup
	var eventErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		requestEvent := createTaskSystemEvent(
			taskId,
			correlationId,
			"API Gateway",
			"api:request:received",
			"Task creation request received",
			marshallJson(params),
		)
		_, err := h.taskSystemEventRepository.Create(requestEvent, 1)
		if err != nil {
			log.Printf("Error creating request event: %v", err)
			eventErr = err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		createdEvent := createTaskSystemEvent(
			taskId,
			correlationId,
			"API Gateway",
			"api:db:task-created",
			"Task created in database",
			"{}",
		)
		_, err := h.taskSystemEventRepository.Create(createdEvent, 2)
		if err != nil {
			log.Printf("Error creating task created event: %v", err)
			eventErr = err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), notificationTimeout)
		defer cancel()

		_, err := h.notificationServiceClient.SendNotification(
			ctx,
			&pb.SendNotificationRequest{
				TaskId:        taskId,
				CorrelationId: correlationId,
				Types:         []pb.NotificationType{0, 1},
			},
			grpc.FailFastCallOption{},
		)

		notificationEvent := createTaskSystemEvent(
			taskId,
			correlationId,
			"API Gateway",
			"api:event:task-created",
			"Task created event emitted",
			"{}",
		)

		_, eventErr := h.taskSystemEventRepository.Create(notificationEvent, 3)
		if eventErr != nil {
			log.Printf("Error creating notification event: %v", eventErr)
		}

		if err != nil {
			log.Printf("Failed to send notification: %v", err)
		} else {
			log.Printf("Notification sent successfully")
		}
	}()

	wg.Wait()

	if eventErr != nil {
		log.Printf("Some events failed to be created: %v", eventErr)
	}

	h.respondWithJSON(w, http.StatusCreated, CreateTaskResponse{
		TaskId: task.ID,
	})
}

// @Summary Update a task
// @Description Updates an existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body UpdateTaskRequest true "Updated task details"
// @Success 200 {object} UpdateTaskResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.respondWithError(w, http.StatusBadRequest, "task ID is required")
		return
	}

	if !commons.IsValidUUID(id) {
		h.respondWithError(w, http.StatusBadRequest, "invalid task ID format")
		return
	}

	task, err := h.taskRepository.GetByID(id)
	if err != nil {
		log.Printf("Error fetching task %s: %v", id, err)
		h.respondWithError(w, http.StatusNotFound, "task not found")
		return
	}

	// Check if user has permission to update the task
	if userID != task.CreatorID && (task.AssigneeID == nil || *task.AssigneeID != userID) {
		h.respondWithError(w, http.StatusForbidden, "You don't have permission to update this task")
		return
	}

	var params UpdateTaskRequest
	if err := h.decodeJSON(r, &params); err != nil {
		log.Printf("Invalid task update request: %v", err)
		h.respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	now := time.Now()
	task.Title = params.Title
	task.Description = params.Description
	task.Status = params.Status
	task.Priority = params.Priority
	task.DueDate = params.DueDate
	task.AssigneeID = params.AssigneeID
	task.UpdatedAt = now

	if err := h.taskRepository.Update(task); err != nil {
		log.Printf("Error updating task %s: %v", id, err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to update task")
		return
	}

	correlationId := uuid.New().String()
	event := createTaskSystemEvent(
		id,
		correlationId,
		"API Gateway",
		"api:event:task-updated",
		"Task updated event emitted",
		"{}",
	)

	_, err = h.taskSystemEventRepository.Create(event, 1)
	if err != nil {
		log.Printf("Error creating task updated event: %v", err)
	}

	h.respondWithJSON(w, http.StatusOK, UpdateTaskResponse{
		TaskId: id,
	})
}

// @Summary Delete a task
// @Description Deletes a task from the system
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.respondWithError(w, http.StatusBadRequest, "task ID is required")
		return
	}

	if !commons.IsValidUUID(id) {
		h.respondWithError(w, http.StatusBadRequest, "invalid task ID format")
		return
	}

	_, err := h.taskRepository.GetByID(id)
	if err != nil {
		log.Printf("Error fetching task %s: %v", id, err)
		h.respondWithError(w, http.StatusNotFound, "task not found")
		return
	}

	if err := h.taskRepository.Delete(id); err != nil {
		log.Printf("Error deleting task %s: %v", id, err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	correlationId := uuid.New().String()
	event := createTaskSystemEvent(
		id,
		correlationId,
		"API Gateway",
		"api:event:task-deleted",
		"Task deleted event emitted",
		"{}",
	)

	_, err = h.taskSystemEventRepository.Create(event, 1)
	if err != nil {
		log.Printf("Error creating task deleted event: %v", err)
	}

	h.respondWithJSON(w, http.StatusOK, map[string]bool{
		"success": true,
	})
}
