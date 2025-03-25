package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	commons "sama/go-task-management/commons"
	pb "sama/go-task-management/commons/api"
	"sama/go-task-management/gateway/config"
	"sama/go-task-management/gateway/interfaces"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type handler struct {
	userRepository              commons.UserRepositoryInterface
	taskRepository              commons.TaskRepositoryInterface
	taskSystemEventRepository   commons.TaskSystemEventRepositoryInterface
	inAppNotificationRepository commons.InAppNotificationRepositoryInterface
	notificationServiceClient   pb.NotificationServiceClient
	passwordResetTokenRepository commons.PasswordResetTokenRepositoryInterface
}

// NewHandler creates a new handler instance
func NewHandler(cfg *config.Config) interfaces.Handler {
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		cfg.NotificationServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	database, err := commons.InitDB()
	if err != nil {
		panic(err)
	}

	return &handler{
		userRepository:              commons.NewPostgresUserRepository(database),
		taskRepository:              commons.NewPostgresTaskRepository(database),
		taskSystemEventRepository:   commons.NewPostgresTaskSystemEventRepository(database),
		inAppNotificationRepository: commons.NewPostgresInAppNotificationRepository(database),
		notificationServiceClient:   pb.NewNotificationServiceClient(conn),
		passwordResetTokenRepository: commons.NewPostgresPasswordResetTokenRepository(database),
	}
}

// Health check endpoint
func (h *handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	log.Println("Registering routes")

	// Health check
	mux.HandleFunc("GET /health", h.Health)

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
