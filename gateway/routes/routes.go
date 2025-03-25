package routes

import (
	"net/http"

	"sama/go-task-management/gateway/interfaces"
	"sama/go-task-management/gateway/middleware"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (r *Router) RegisterRoutes(handler interfaces.Handler, corsConfig middleware.CorsConfig, authConfig middleware.AuthConfig) {
	// Health check
	r.mux.HandleFunc("GET /health", handler.Health)

	// Auth routes
	r.mux.HandleFunc("POST /api/v1/auth/signin", handler.Signin)
	r.mux.HandleFunc("POST /api/v1/auth/signup", handler.Signup)
	r.mux.HandleFunc("POST /api/v1/auth/refresh-token", handler.RefreshToken)
	r.mux.HandleFunc("POST /api/v1/auth/signout", handler.Signout)
	r.mux.HandleFunc("POST /api/v1/auth/forgot-password", handler.ForgotPassword)
	r.mux.HandleFunc("POST /api/v1/auth/reset-password", handler.ResetPassword)

	// Task routes
	r.mux.HandleFunc("GET /api/v1/tasks/{id}", handler.GetTask)
	r.mux.HandleFunc("GET /api/v1/tasks", handler.GetAllTasks)
	r.mux.HandleFunc("POST /api/v1/tasks", handler.CreateTask)
	r.mux.HandleFunc("PUT /api/v1/tasks/{id}", handler.UpdateTask)
	r.mux.HandleFunc("DELETE /api/v1/tasks/{id}", handler.DeleteTask)

	// Notification routes
	r.mux.HandleFunc("GET /api/v1/notifications", handler.GetAllInAppNotifications)
	r.mux.HandleFunc("POST /api/v1/notifications/{id}/read", handler.UpdateOnRead)
	r.mux.HandleFunc("DELETE /api/v1/notifications/{id}", handler.DeleteInAppNotification)

	// System event routes
	r.mux.HandleFunc("GET /api/v1/task-system-events", handler.GetAllTaskSystemEvents)
}

func (r *Router) GetMux() *http.ServeMux {
	return r.mux
}
