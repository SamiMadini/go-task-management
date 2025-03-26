package handlers

import "net/http"

type HandlerWrapper struct {
	Base              *BaseHandler
	Auth              *AuthHandler
	Task              *TaskHandler
	InAppNotification *InAppNotificationHandler
	TaskSystemEvent   *TaskSystemEventHandler
}

func (h *HandlerWrapper) Health(w http.ResponseWriter, r *http.Request) {
	h.Base.Health(w, r)
}

func (h *HandlerWrapper) SignIn(w http.ResponseWriter, r *http.Request) {
	h.Auth.SignIn(w, r)
}

func (h *HandlerWrapper) SignUp(w http.ResponseWriter, r *http.Request) {
	h.Auth.SignUp(w, r)
}

func (h *HandlerWrapper) SignOut(w http.ResponseWriter, r *http.Request) {
	h.Auth.SignOut(w, r)
}

func (h *HandlerWrapper) RefreshToken(w http.ResponseWriter, r *http.Request) {
	h.Auth.RefreshToken(w, r)
}

func (h *HandlerWrapper) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	h.Auth.ForgotPassword(w, r)
}

func (h *HandlerWrapper) ResetPassword(w http.ResponseWriter, r *http.Request) {
	h.Auth.ResetPassword(w, r)
}

func (h *HandlerWrapper) GetTask(w http.ResponseWriter, r *http.Request) {
	h.Task.GetTask(w, r)
}

func (h *HandlerWrapper) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	h.Task.GetAllTasks(w, r)
}

func (h *HandlerWrapper) CreateTask(w http.ResponseWriter, r *http.Request) {
	h.Task.CreateTask(w, r)
}

func (h *HandlerWrapper) UpdateTask(w http.ResponseWriter, r *http.Request) {
	h.Task.UpdateTask(w, r)
}

func (h *HandlerWrapper) DeleteTask(w http.ResponseWriter, r *http.Request) {
	h.Task.DeleteTask(w, r)
}

func (h *HandlerWrapper) GetAllInAppNotifications(w http.ResponseWriter, r *http.Request) {
	h.InAppNotification.GetUserNotifications(w, r)
}

func (h *HandlerWrapper) UpdateOnRead(w http.ResponseWriter, r *http.Request) {
	h.InAppNotification.MarkAsRead(w, r)
}

func (h *HandlerWrapper) DeleteInAppNotification(w http.ResponseWriter, r *http.Request) {
	h.InAppNotification.DeleteNotification(w, r)
}

func (h *HandlerWrapper) GetAllTaskSystemEvents(w http.ResponseWriter, r *http.Request) {
	h.TaskSystemEvent.GetAllTaskSystemEvents(w, r)
}
