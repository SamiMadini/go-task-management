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
	passwordResetTokenRepository commons.PasswordResetTokenRepositoryInterface
}

func NewHandler(
	userRepository commons.UserRepositoryInterface,
	taskRepository commons.TaskRepositoryInterface,
	taskSystemEventRepository commons.TaskSystemEventRepositoryInterface,
	inAppNotificationRepository commons.InAppNotificationRepositoryInterface,
	notificationServiceClient pb.NotificationServiceClient,
	passwordResetTokenRepository commons.PasswordResetTokenRepositoryInterface,
) *handler {
	return &handler{
		userRepository:             userRepository,
		taskRepository:            taskRepository,
		taskSystemEventRepository:  taskSystemEventRepository,
		inAppNotificationRepository: inAppNotificationRepository,
		notificationServiceClient:   notificationServiceClient,
		passwordResetTokenRepository: passwordResetTokenRepository,
	}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	log.Println("Registering routes")

	// Health check
	mux.HandleFunc("GET /health", h.health)

	// Auth routes
	mux.HandleFunc("POST /api/v1/auth/signin", h.Signin)
	mux.HandleFunc("POST /api/v1/auth/signup", h.Signup)
	mux.HandleFunc("POST /api/v1/auth/refresh-token", h.RefreshToken)
	mux.HandleFunc("POST /api/v1/auth/signout", h.Signout)
	mux.HandleFunc("POST /api/v1/auth/forgot-password", h.ForgotPassword)
	mux.HandleFunc("POST /api/v1/auth/reset-password", h.ResetPassword)

	// Task endpoints
	mux.HandleFunc("GET /api/v1/tasks/{id}", h.GetTask)
	mux.HandleFunc("GET /api/v1/tasks", h.GetAllTasks)
	mux.HandleFunc("POST /api/v1/tasks", h.CreateTask)
	mux.HandleFunc("PUT /api/v1/tasks/{id}", h.UpdateTask)
	mux.HandleFunc("DELETE /api/v1/tasks/{id}", h.DeleteTask)

	// Notification endpoints
	mux.HandleFunc("GET /api/v1/notifications", h.GetAllInAppNotifications)
	mux.HandleFunc("POST /api/v1/notifications/{id}/read", h.UpdateOnRead)
	mux.HandleFunc("DELETE /api/v1/notifications/{id}", h.DeleteInAppNotification)

	// System event endpoints
	mux.HandleFunc("GET /api/v1/task-system-events", h.GetAllTaskSystemEvents)
}

func (h *handler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
