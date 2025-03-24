package main

import (
	"log"
	"net/http"

	commons "sama/go-task-management/commons"
	pb "sama/go-task-management/commons/api"
)

type handler struct {
	taskRepository              commons.TaskRepositoryInterface
	inAppNotificationRepository commons.InAppNotificationRepositoryInterface
	taskSystemEventRepository   commons.TaskSystemEventRepositoryInterface
	notificationServiceClient   pb.NotificationServiceClient
	userRepository             commons.UserRepositoryInterface
}

func NewHandler(
	userRepository commons.UserRepositoryInterface,
	taskRepository commons.TaskRepositoryInterface,
	taskSystemEventRepository commons.TaskSystemEventRepositoryInterface,
	inAppNotificationRepository commons.InAppNotificationRepositoryInterface,
	notificationServiceClient pb.NotificationServiceClient,
) *handler {
	return &handler{
		taskRepository:              taskRepository,
		inAppNotificationRepository: inAppNotificationRepository,
		taskSystemEventRepository:   taskSystemEventRepository,
		notificationServiceClient:   notificationServiceClient,
		userRepository:             userRepository,
	}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	log.Println("Registering routes")

	// Health check
	mux.HandleFunc("GET /api/_health", h.health)

	// Auth endpoints (no authentication required)
	mux.HandleFunc("POST /api/auth/signup", h.Signup)
	mux.HandleFunc("POST /api/auth/signin", h.Signin)
	mux.HandleFunc("POST /api/auth/refresh", h.RefreshToken)
	mux.HandleFunc("POST /api/auth/signout", h.Signout)

	// Task endpoints
	mux.HandleFunc("GET /api/tasks/{id}", h.GetTask)
	mux.HandleFunc("GET /api/tasks", h.GetAllTasks)
	mux.HandleFunc("POST /api/tasks", h.CreateTask)
	mux.HandleFunc("PUT /api/tasks/{id}", h.UpdateTask)
	mux.HandleFunc("DELETE /api/tasks/{id}", h.DeleteTask)

	// Notification endpoints
	mux.HandleFunc("GET /api/notifications", h.GetAllInAppNotifications)
	mux.HandleFunc("POST /api/notifications/{id}/read", h.UpdateOnRead)
	mux.HandleFunc("DELETE /api/notifications/{id}", h.DeleteInAppNotification)

	// System event endpoints
	mux.HandleFunc("GET /api/task-system-events", h.GetAllTaskSystemEvents)
}

func (h *handler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
