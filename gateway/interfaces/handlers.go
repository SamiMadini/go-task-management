package interfaces

import "net/http"

type HealthHandler interface {
	Health(w http.ResponseWriter, r *http.Request)
}

type AuthHandler interface {
	Signin(w http.ResponseWriter, r *http.Request)
	Signup(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	Signout(w http.ResponseWriter, r *http.Request)
	ForgotPassword(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)
}

type TaskHandler interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	GetAllTasks(w http.ResponseWriter, r *http.Request)
	CreateTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
}

type NotificationHandler interface {
	GetAllInAppNotifications(w http.ResponseWriter, r *http.Request)
	UpdateOnRead(w http.ResponseWriter, r *http.Request)
	DeleteInAppNotification(w http.ResponseWriter, r *http.Request)
}

type SystemEventHandler interface {
	GetAllTaskSystemEvents(w http.ResponseWriter, r *http.Request)
}

type Handler interface {
	HealthHandler
	AuthHandler
	TaskHandler
	NotificationHandler
	SystemEventHandler
}
