package main

import (
	"errors"
	"time"
)

type GetTaskResponse struct {
	ID          string                    `json:"id"`
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	Status      string                    `json:"status"`
	Priority    int                       `json:"priority"`
	DueDate     time.Time                 `json:"due_date"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
	Events      []TaskSystemEventResponse `json:"events"`
}

type GetAllTasksResponse struct {
	Tasks []GetTaskResponse `json:"tasks"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"max=500"`
	Status      string `json:"status" validate:"required,oneof=todo in_progress done"`
	Priority    int    `json:"priority" validate:"required,oneof=1 2 3"`
	DueDate     string `json:"due_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

func (r *CreateTaskRequest) Validate() error {
	if r.Title == "" {
		return errors.New("title is required")
	}
	if len(r.Title) < 3 {
		return errors.New("title must be at least 3 characters long")
	}
	if len(r.Title) > 100 {
		return errors.New("title must be less than 100 characters")
	}
	if len(r.Description) > 500 {
		return errors.New("description must be less than 500 characters")
	}
	if r.Status == "" {
		return errors.New("status is required")
	}
	if r.Status != "todo" && r.Status != "in_progress" && r.Status != "done" {
		return errors.New("status must be one of: todo, in_progress, done")
	}
	if r.Priority < 1 || r.Priority > 3 {
		return errors.New("priority must be between 1 and 3")
	}
	if r.DueDate != "" {
		_, err := time.Parse(time.RFC3339, r.DueDate)
		if err != nil {
			return errors.New("due_date must be in RFC3339 format (e.g., 2024-03-19T15:04:05Z07:00)")
		}
	}
	return nil
}

type CreateTaskResponse struct {
	TaskId string `json:"task_id"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title" validate:"omitempty,min=3,max=100"`
	Description string `json:"description" validate:"omitempty,max=500"`
	Status      string `json:"status" validate:"omitempty,oneof=todo in_progress done"`
	Priority    int    `json:"priority" validate:"omitempty,oneof=1 2 3"`
	DueDate     string `json:"due_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

func (r *UpdateTaskRequest) Validate() error {
	if r.Title != "" {
		if len(r.Title) < 3 {
			return errors.New("title must be at least 3 characters long")
		}
		if len(r.Title) > 100 {
			return errors.New("title must be less than 100 characters")
		}
	}
	if r.Description != "" && len(r.Description) > 500 {
		return errors.New("description must be less than 500 characters")
	}
	if r.Status != "" && r.Status != "todo" && r.Status != "in_progress" && r.Status != "done" {
		return errors.New("status must be one of: todo, in_progress, done")
	}
	if r.Priority != 0 && (r.Priority < 1 || r.Priority > 3) {
		return errors.New("priority must be between 1 and 3")
	}
	if r.DueDate != "" {
		_, err := time.Parse(time.RFC3339, r.DueDate)
		if err != nil {
			return errors.New("due_date must be in RFC3339 format (e.g., 2024-03-19T15:04:05Z07:00)")
		}
	}
	return nil
}

type UpdateTaskResponse struct {
	TaskId string `json:"task_id"`
}

type GetAllInAppNotificationsResponse struct {
	InAppNotifications []InAppNotificationResponse `json:"in_app_notifications"`
}

type InAppNotificationResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsRead      bool      `json:"is_read"`
	ReadAt      time.Time `json:"read_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateOnReadRequest struct {
	IsRead bool `json:"is_read" validate:"required"`
}

func (r *UpdateOnReadRequest) Validate() error {
	return nil
}

type UpdateOnReadResponse struct {
	Success bool `json:"success"`
}

type TaskSystemEventResponse struct {
	ID            string    `json:"id"`
	TaskId        string    `json:"task_id"`
	CorrelationId string    `json:"correlation_id"`
	Origin        string    `json:"origin"`
	Action        string    `json:"action"`
	Message       string    `json:"message"`
	JsonData      string    `json:"json_data"`
	EmitAt        time.Time `json:"emit_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type GetAllTaskSystemEventsResponse struct {
	Events []TaskSystemEventResponse `json:"events"`
}
