package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"

	commons "sama/go-task-management/commons"
	pb "sama/go-task-management/commons/api"

	"google.golang.org/grpc"
)

// @Summary Get a task by ID
// @Description Retrieves a specific task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} GetTaskResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id} [get]
func (h *handler) GetTask(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    if id == "" {
        commons.WriteJSONError(w, http.StatusBadRequest, "task ID is required")
        return
    }

    if !commons.IsValidUUID(id) {
        commons.WriteJSONError(w, http.StatusBadRequest, "invalid task ID format")
        return
    }

    task, err := h.taskRepository.GetByID(id)
    if err != nil {
        log.Printf("Error fetching task %s: %v", id, err)
        commons.InternalServerErrorHandler(w)
        return
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
        Events:      events,
    }

    commons.WriteJSON(w, http.StatusOK, response)
}

// @Summary Get all tasks
// @Description Retrieves all tasks in the system
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} GetAllTasksResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks [get]
func (h *handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
    tasks, err := h.taskRepository.GetAll()
    if err != nil {
        log.Printf("Failed to get all tasks: %v", err)
        commons.InternalServerErrorHandler(w)
        return
    }

    if len(tasks) == 0 {
        commons.WriteJSON(w, http.StatusOK, GetAllTasksResponse{
            Tasks: []GetTaskResponse{},
        })
        return
    }

    response := GetAllTasksResponse{
        Tasks: make([]GetTaskResponse, len(tasks)),
    }

    for i, task := range tasks {
        events := make([]TaskSystemEventResponse, len(task.Events))
        for j, event := range task.Events {
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
            ID:          task.ID,
            Title:       task.Title,
            Description: task.Description,
            Status:      task.Status,
            Priority:    task.Priority,
            DueDate:     task.DueDate,
            CreatedAt:   task.CreatedAt,
            UpdatedAt:   task.UpdatedAt,
            Events:      events,
        }
    }

    commons.WriteJSON(w, http.StatusOK, response)
}

// @Summary Create a new task
// @Description Creates a new task in the system
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body CreateTaskRequest true "Task details"
// @Success 201 {object} CreateTaskResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks [post]
func (h *handler) CreateTask(w http.ResponseWriter, r *http.Request) {
    var params CreateTaskRequest
    if err := commons.ReadJSON(r, &params); err != nil {
        log.Printf("Invalid task creation request: %v", err)
        commons.WriteJSONError(w, http.StatusBadRequest, "Invalid request format")
        return
    }

    if err := params.Validate(); err != nil {
        log.Printf("Validation error: %v", err)
        commons.WriteJSONError(w, http.StatusBadRequest, err.Error())
        return
    }

    now := time.Now()
    taskId := uuid.New().String()
    correlationId := uuid.New().String()

    task := commons.Task{
        ID:          taskId,
        Title:       params.Title,
        Description: params.Description,
        Status:      params.Status,
        Priority:    params.Priority,
        DueDate:     parseDueDate(params.DueDate, now),
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    var result commons.Task
    result, err := h.taskRepository.Create(task)
    if err != nil {
        log.Printf("Error creating task: %v", err)
        commons.InternalServerErrorHandler(w)
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
            eventErr = err
        }

        if err != nil {
            log.Printf("Failed to send notification: %v", err)
            eventErr = err
        } else {
            log.Printf("Notification sent successfully")
        }
    }()

    wg.Wait()

    if eventErr != nil {
        log.Printf("Some events failed to be created: %v", eventErr)
    }

    commons.WriteJSON(w, http.StatusCreated, CreateTaskResponse{
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
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id} [put]
func (h *handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    if id == "" {
        commons.WriteJSONError(w, http.StatusBadRequest, "task ID is required")
        return
    }

    if !commons.IsValidUUID(id) {
        commons.WriteJSONError(w, http.StatusBadRequest, "invalid task ID format")
        return
    }

    _, err := h.taskRepository.GetByID(id)
    if err != nil {
        log.Printf("Error fetching task %s: %v", id, err)
        commons.WriteJSONError(w, http.StatusNotFound, "task not found")
        return
    }

    var params UpdateTaskRequest
    if err := commons.ReadJSON(r, &params); err != nil {
        log.Printf("Invalid task update request: %v", err)
        commons.WriteJSONError(w, http.StatusBadRequest, "Invalid request format")
        return
    }

    if err := params.Validate(); err != nil {
        log.Printf("Validation error: %v", err)
        commons.WriteJSONError(w, http.StatusBadRequest, err.Error())
        return
    }

    now := time.Now()
    task := commons.Task{
        ID:          id,
        Title:       params.Title,
        Description: params.Description,
        Status:      params.Status,
        Priority:    params.Priority,
        DueDate:     parseDueDate(params.DueDate, time.Time{}),
        UpdatedAt:   now,
    }

    if err := h.taskRepository.Update(task); err != nil {
        log.Printf("Error updating task %s: %v", id, err)
        commons.InternalServerErrorHandler(w)
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

    commons.WriteJSON(w, http.StatusOK, UpdateTaskResponse{
        TaskId: id,
    })
}

// @Summary Delete a task
// @Description Deletes a task from the system
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id} [delete]
func (h *handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    if id == "" {
        commons.WriteJSONError(w, http.StatusBadRequest, "task ID is required")
        return
    }

    if !commons.IsValidUUID(id) {
        commons.WriteJSONError(w, http.StatusBadRequest, "invalid task ID format")
        return
    }

    _, err := h.taskRepository.GetByID(id)
    if err != nil {
        log.Printf("Error fetching task %s: %v", id, err)
        commons.WriteJSONError(w, http.StatusNotFound, "task not found")
        return
    }

    if err := h.taskRepository.Delete(id); err != nil {
        log.Printf("Error deleting task %s: %v", id, err)
        commons.InternalServerErrorHandler(w)
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

    commons.WriteJSON(w, http.StatusOK, map[string]bool{
        "success": true,
    })
}
