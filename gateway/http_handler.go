package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"

	commons "sama/go-task-management/commons"
	pb "sama/go-task-management/commons/api"

	"google.golang.org/grpc"
)

type handler struct {
	taskRepository              commons.TaskRepositoryInterface
	inAppNotificationRepository commons.InAppNotificationRepositoryInterface
	taskSystemEventRepository   commons.TaskSystemEventRepositoryInterface
	notificationServiceClient   pb.NotificationServiceClient
}

func NewHandler(
	taskRepository commons.TaskRepositoryInterface,
	inAppNotificationRepository commons.InAppNotificationRepositoryInterface,
	taskSystemEventRepository commons.TaskSystemEventRepositoryInterface,
	notificationServiceClient pb.NotificationServiceClient,
) *handler {
	return &handler{
		taskRepository:              taskRepository,
		inAppNotificationRepository: inAppNotificationRepository,
		taskSystemEventRepository:   taskSystemEventRepository,
		notificationServiceClient:   notificationServiceClient,
	}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	log.Println("Registering routes")

	mux.HandleFunc("GET /api/_health", h.health)

	mux.HandleFunc("GET /api/tasks/{id}", h.GetTask)
	mux.HandleFunc("GET /api/tasks", h.GetAllTasks)
	mux.HandleFunc("POST /api/tasks", h.CreateTask)
	mux.HandleFunc("PUT /api/tasks/{id}", h.UpdateTask)
	mux.HandleFunc("DELETE /api/tasks/{id}", h.DeleteTask)

	mux.HandleFunc("GET /api/notifications", h.GetAllInAppNotifications)
	mux.HandleFunc("POST /api/notifications/{id}/read", h.UpdateOnRead)
	mux.HandleFunc("DELETE /api/notifications/{id}", h.DeleteInAppNotification)

	mux.HandleFunc("GET /api/task-system-events", h.GetAllTaskSystemEvents)
}

func (h *handler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	task, err := h.taskRepository.GetByID(id)
	if err != nil {
		commons.InternalServerErrorHandler(w)
		return
	}

	commons.WriteJSON(w, http.StatusOK, task)
}

// GetAllTasks godoc
// @Summary Get all tasks
// @Description Get all tasks from the database
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {array} models.Task
// @Router /tasks [get]
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
	err := commons.ReadJSON(r, &params)

	if err != nil {
		commons.InternalServerErrorHandler(w)
		return
	}

	now := time.Now()
	taskId := uuid.New().String()
	correlationId := uuid.New().String()

	var wg = &sync.WaitGroup{}
	
	paramsTask := commons.Task{
		ID:          taskId,
		Title:       params.Title,
		Description: params.Description,
		Status:      params.Status,
		Priority:    params.Priority,
		DueDate:     func() time.Time {
			if params.DueDate == "" {
				return now
			}
			t, err := time.Parse(time.RFC3339, params.DueDate)
			if err != nil {
				log.Printf("Error parsing due date: %v", err)
				return now
			}
			return t
		}(),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	var task commons.Task
	wg.Add(1)
	go func() {
		defer wg.Done()
		task, err = h.taskRepository.Create(paramsTask)
		if err != nil {
			commons.InternalServerErrorHandler(w)
			return
		}
	}()

	eventTaskParams := commons.TaskSystemEvent{
		TaskId: taskId,
		CorrelationId: correlationId,
		Origin: "API Gateway",
		Action: "api:request:received",
		Message: "Task creation request received",
		JsonData: func() string {
			jsonData, err := json.Marshal(params)
			if err != nil {
				return "{}"
			}
			return string(jsonData)
		}(),
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		h.taskSystemEventRepository.Create(eventTaskParams, 1)
	}()

	eventTaskCreated := commons.TaskSystemEvent{
		TaskId: taskId,
		CorrelationId: correlationId,
		Origin: "API Gateway",
		Action: "api:db:task-created",
		Message: "Task created in database",
		JsonData: "{}",
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		h.taskSystemEventRepository.Create(eventTaskCreated, 2)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		
		_, err := h.notificationServiceClient.SendNotification(ctx, &pb.SendNotificationRequest{
			TaskId: taskId,
			CorrelationId: correlationId,
			Types:  []pb.NotificationType{0, 1},
		}, grpc.FailFastCallOption{})

		eventTaskNotificationSent := commons.TaskSystemEvent{
			TaskId: taskId,
			CorrelationId: correlationId,
			Origin: "API Gateway",
			Action: "api:event:task-created",
			Message: "Task created event emitted",
			JsonData: "{}",
		}
	
		h.taskSystemEventRepository.Create(eventTaskNotificationSent, 3)

		if err != nil {
			log.Printf("Failed to send event: %v", err)
		} else {
			log.Printf("Event sent successfully")
		}
	}()

	wg.Wait()

	commons.WriteJSON(w, http.StatusOK, task)
}

func (h *handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var params UpdateTaskRequest
	err := commons.ReadJSON(r, &params)

	if err != nil {
		commons.InternalServerErrorHandler(w)
		return
	}

	paramsTask := commons.Task{
		ID:          id,
		Title:       params.Title,
		Description: params.Description,
		Status:      params.Status,
		Priority:    params.Priority,
		DueDate:     func() time.Time {
			if params.DueDate == "" {
				return time.Time{} // Zero time if empty
			}
			t, err := time.Parse(time.RFC3339, params.DueDate)
			if err != nil {
				log.Printf("Error parsing due date: %v", err)
				return time.Time{} // Return zero time on parsing error
			}
			return t
		}(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = h.taskRepository.Update(paramsTask)
	if err != nil {
		commons.InternalServerErrorHandler(w)
	}

	correlationId := uuid.New().String()
	eventTaskUpdated := commons.TaskSystemEvent{
		TaskId: id,
		CorrelationId: correlationId,
		Origin: "API Gateway",
		Action: "api:event:task-updated",
		Message: "Task updated event emitted",
		JsonData: "{}",
	}

	h.taskSystemEventRepository.Create(eventTaskUpdated, 1)

	commons.WriteJSON(w, http.StatusOK, "OK")
}

func (h *handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.taskRepository.Delete(id)
	if err != nil {
		commons.InternalServerErrorHandler(w)
	}

	correlationId := uuid.New().String()
	eventTaskDeleted := commons.TaskSystemEvent{
		TaskId: id,
		CorrelationId: correlationId,
		Origin: "API Gateway",
		Action: "api:event:task-deleted",
		Message: "Task deleted event emitted",
		JsonData: "{}",
	}

	h.taskSystemEventRepository.Create(eventTaskDeleted, 1)

	commons.WriteJSON(w, http.StatusOK, "OK")
}

// InAppNotifications

func (h *handler) GetAllInAppNotifications(w http.ResponseWriter, r *http.Request) {
	inAppNotifications, err := h.inAppNotificationRepository.GetAll()
	if err != nil {
		commons.InternalServerErrorHandler(w)
		return
	}

	if len(inAppNotifications) == 0 {
		commons.WriteJSON(w, http.StatusOK, []commons.InAppNotification{})
		return
	}

	commons.WriteJSON(w, http.StatusOK, inAppNotifications)
}

func (h *handler) UpdateOnRead(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var params UpdateOnReadRequest
	err := commons.ReadJSON(r, &params)

	if err != nil {
		commons.InternalServerErrorHandler(w)
		return
	}

	log.Println(params)

	err = h.inAppNotificationRepository.UpdateOnRead(id, params.IsRead)
	if err != nil {
		commons.InternalServerErrorHandler(w)
	}

	commons.WriteJSON(w, http.StatusOK, "OK")
}

func (h *handler) DeleteInAppNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.inAppNotificationRepository.Delete(id)
	if err != nil {
		commons.InternalServerErrorHandler(w)
	}

	commons.WriteJSON(w, http.StatusOK, "OK")
}


// TaskSystemEvents

func (h *handler) GetAllTaskSystemEvents(w http.ResponseWriter, r *http.Request) {
	taskSystemEvents, err := h.taskSystemEventRepository.GetAll()
	if err != nil {
		commons.InternalServerErrorHandler(w)
		return
	}

	if len(taskSystemEvents) == 0 {
		commons.WriteJSON(w, http.StatusOK, []commons.TaskSystemEvent{})
		return
	}

	commons.WriteJSON(w, http.StatusOK, taskSystemEvents)
}
