package routes

import (
	"net/http"
	"sama/go-task-management/gateway/interfaces"
	"sama/go-task-management/gateway/middleware"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	router chi.Router
}

func NewRouter() *Router {
	return &Router{
		router: chi.NewRouter(),
	}
}

func (r *Router) RegisterRoutes(handler interfaces.Handler, authConfig middleware.AuthConfig) {
	// Health check
	r.router.Get("/health", handler.Health)

	// Swagger UI routes
	fs := http.FileServer(http.Dir("./gateway/docs"))
	r.router.Handle("/swagger/*", http.StripPrefix("/swagger/", fs))

	// Public routes
	r.router.Group(func(router chi.Router) {
		// Auth routes
		router.Post("/api/v1/auth/signin", handler.SignIn)
		router.Post("/api/v1/auth/signup", handler.SignUp)
		router.Post("/api/v1/auth/refresh", handler.RefreshToken)
		router.Post("/api/v1/auth/forgot-password", handler.ForgotPassword)
		router.Post("/api/v1/auth/reset-password", handler.ResetPassword)
	})

	// Protected routes
	r.router.Group(func(router chi.Router) {
		router.Use(middleware.AuthMiddleware(authConfig))

		// Auth routes
		router.Post("/api/v1/auth/signout", handler.SignOut)

		// Task routes
		router.Get("/api/v1/tasks/{id}", handler.GetTask)
		router.Get("/api/v1/tasks", handler.GetAllTasks)
		router.Post("/api/v1/tasks", handler.CreateTask)
		router.Put("/api/v1/tasks/{id}", handler.UpdateTask)
		router.Delete("/api/v1/tasks/{id}", handler.DeleteTask)

		// Notification routes
		router.Get("/api/v1/notifications", handler.GetAllInAppNotifications)
		router.Post("/api/v1/notifications/{id}/read", handler.UpdateOnRead)
		router.Delete("/api/v1/notifications/{id}", handler.DeleteInAppNotification)

		// System event routes
		router.Get("/api/v1/task-system-events", handler.GetAllTaskSystemEvents)
	})
}

func (r *Router) GetRouter() chi.Router {
	return r.router
}
