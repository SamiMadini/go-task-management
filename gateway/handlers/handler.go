package handlers

import (
	"net/http"

	"sama/go-task-management/gateway/config"
	"sama/go-task-management/gateway/interfaces"
)

type Handler struct {
	*BaseHandler
	*AuthHandler
	*TaskHandler
	*NotificationHandler
	*SystemEventHandler
}

func NewHandler(cfg *config.Config) (interfaces.Handler, error) {
	base, err := NewBaseHandler(cfg)
	if err != nil {
		return nil, err
	}

	return &Handler{
		BaseHandler:         base,
		AuthHandler:         NewAuthHandler(base),
		TaskHandler:         NewTaskHandler(base),
		NotificationHandler: NewNotificationHandler(base),
		SystemEventHandler:  NewSystemEventHandler(base),
	}, nil
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	h.BaseHandler.Health(w, r)
}

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	h.AuthHandler.Signin(w, r)
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	h.AuthHandler.Signup(w, r)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	h.AuthHandler.RefreshToken(w, r)
}

func (h *Handler) Signout(w http.ResponseWriter, r *http.Request) {
	h.AuthHandler.Signout(w, r)
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	h.AuthHandler.ForgotPassword(w, r)
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	h.AuthHandler.ResetPassword(w, r)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	h.TaskHandler.GetTask(w, r)
}

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	h.TaskHandler.GetAllTasks(w, r)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	h.TaskHandler.CreateTask(w, r)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	h.TaskHandler.UpdateTask(w, r)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	h.TaskHandler.DeleteTask(w, r)
}

func (h *Handler) GetAllInAppNotifications(w http.ResponseWriter, r *http.Request) {
	h.NotificationHandler.GetAllInAppNotifications(w, r)
}

func (h *Handler) UpdateOnRead(w http.ResponseWriter, r *http.Request) {
	h.NotificationHandler.UpdateOnRead(w, r)
}

func (h *Handler) DeleteInAppNotification(w http.ResponseWriter, r *http.Request) {
	h.NotificationHandler.DeleteInAppNotification(w, r)
}

func (h *Handler) GetAllTaskSystemEvents(w http.ResponseWriter, r *http.Request) {
	h.SystemEventHandler.GetAllTaskSystemEvents(w, r)
}
