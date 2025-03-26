package in_app_notification

import (
	"sama/go-task-management/commons"
)

type CreateNotificationInput struct {
	UserID  string `json:"user_id"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

type NotificationResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Message     string `json:"message"`
	IsRead      bool   `json:"is_read"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func toNotificationResponse(notification commons.InAppNotification) NotificationResponse {
	return NotificationResponse{
		ID:        notification.ID,
		UserID:    notification.UserID,
		Title:     notification.Title,
		Message:   notification.Description,
		IsRead:    notification.IsRead,
		CreatedAt: notification.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: notification.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
