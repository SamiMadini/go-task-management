package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	commons "sama/go-task-management/commons"
	pb "sama/go-task-management/commons/api"

	"google.golang.org/grpc"
)

type handler struct {
	taskRepository              commons.TaskRepositoryInterface
	inAppNotificationRepository commons.InAppNotificationRepositoryInterface
	notificationServiceClient   pb.NotificationServiceClient
}

func NewHandler(
	taskRepository commons.TaskRepositoryInterface,
	inAppNotificationRepository commons.InAppNotificationRepositoryInterface,
	notificationServiceClient pb.NotificationServiceClient,
) *handler {
	return &handler{
		taskRepository:              taskRepository,
		inAppNotificationRepository: inAppNotificationRepository,
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
	log.Println("GetAllTasks")

	tasks, err := h.taskRepository.GetAll()
	if err != nil {
		log.Printf("Failed to get all tasks: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	log.Println("Tasks: ", tasks)

	if len(tasks) == 0 {
		commons.WriteJSON(w, http.StatusOK, []commons.Task{})
		return
	}

	log.Println("Writing JSON")
	commons.WriteJSON(w, http.StatusOK, tasks)
}

func (h *handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("createTask")

	var params CreateTaskRequest
	err := commons.ReadJSON(r, &params)

	if err != nil {
		commons.InternalServerErrorHandler(w)
		return
	}

	log.Println(params)
	
	paramsTask := commons.Task{
		Title:       params.Title,
		Description: params.Description,
		Status:      params.Status,
		Priority:    params.Priority,
		DueDate:     time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	task, err := h.taskRepository.Create(paramsTask)
	if err != nil {
		commons.InternalServerErrorHandler(w)
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		
		_, err := h.notificationServiceClient.SendNotification(ctx, &pb.SendNotificationRequest{
			TaskId: task.ID,
			Types:  []pb.NotificationType{0, 1},
		}, grpc.FailFastCallOption{})
		if err != nil {
			log.Printf("Failed to send notification: %v", err)
		} else {
			log.Printf("Notification sent successfully")
		}
	}()

	commons.WriteJSON(w, http.StatusOK, task)
}

func (h *handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	fmt.Println("updateTask", id)

	var params UpdateTaskRequest
	err := commons.ReadJSON(r, &params)

	if err != nil {
		commons.InternalServerErrorHandler(w)
		return
	}

	log.Println(params)

	paramsTask := commons.Task{
		ID:          id,
		Title:       params.Title,
		Description: params.Description,
		Status:      params.Status,
		Priority:    params.Priority,
		DueDate:     time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}


	err = h.taskRepository.Update(paramsTask)
	if err != nil {
		commons.InternalServerErrorHandler(w)
	}

	commons.WriteJSON(w, http.StatusOK, "OK")
}

func (h *handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.taskRepository.Delete(id)
	if err != nil {
		commons.InternalServerErrorHandler(w)
	}

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

	log.Println("UpdateOnRead", id)

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
