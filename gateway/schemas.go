package main

import (
	"time"
)

type GetTaskResponse struct {
	TaskId string
	Title string
	Description string
	Status string
	Priority int
	DueDate time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateTaskRequest struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Status string `json:"status"`
	Priority int `json:"priority"`
	DueDate string `json:"due_date"`
}

type CreateTaskResponse struct {
	TaskId string
}

type UpdateTaskRequest struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Status string `json:"status"`
	Priority int `json:"priority"`
	DueDate string `json:"due_date"`
}	

type UpdateTaskResponse struct {
	TaskId string
}

type GetAllInAppNotificationsResponse struct {
	InAppNotifications []InAppNotificationResponse
}

type InAppNotificationResponse struct {
	ID string
	Title string
	Description string
	IsRead bool
	ReadAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateOnReadRequest struct {
	IsRead bool
}

type UpdateOnReadResponse struct {
	Success bool
}
