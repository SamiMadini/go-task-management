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

func (h *handler) GetTask(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    if id == "" {
        commons.WriteJSONError(w, http.StatusBadRequest, "task ID is required")
        return
    }

    task, err := h.taskRepository.GetByID(id)
    if err != nil {
        log.Printf("Error fetching task %s: %v", id, err)
        commons.InternalServerErrorHandler(w)
        return
    }

    commons.WriteJSON(w, http.StatusOK, task)
}

func (h *handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
    tasks, err := h.taskRepository.GetAll()
    if err != nil {
        log.Printf("Failed to get all tasks: %v", err)
        commons.InternalServerErrorHandler(w)
        return
    }

    if len(tasks) == 0 {
        commons.WriteJSON(w, http.StatusOK, []commons.Task{})
        return
    }

    commons.WriteJSON(w, http.StatusOK, tasks)
}

func (h *handler) CreateTask(w http.ResponseWriter, r *http.Request) {
    var params CreateTaskRequest
    if err := commons.ReadJSON(r, &params); err != nil {
        log.Printf("Invalid task creation request: %v", err)
        commons.WriteJSONError(w, http.StatusBadRequest, "Invalid request format")
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

    var wg sync.WaitGroup
    var createErr error

    wg.Add(1)
    go func() {
        defer wg.Done()
        var result commons.Task
        result, createErr = h.taskRepository.Create(task)
        if createErr != nil {
            log.Printf("Error creating task: %v", createErr)
            return
        }
        task = result
    }()

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

    if createErr != nil {
        commons.InternalServerErrorHandler(w)
        return
    }

    commons.WriteJSON(w, http.StatusCreated, task)
}

func (h *handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    if id == "" {
        commons.WriteJSONError(w, http.StatusBadRequest, "task ID is required")
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

    commons.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (h *handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    if id == "" {
        commons.WriteJSONError(w, http.StatusBadRequest, "task ID is required")
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

    commons.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
