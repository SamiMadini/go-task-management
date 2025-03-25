package handlers

import (
	"time"
)

type SignupResponse struct {
	User struct {
		ID     string `json:"id"`
		Handle string `json:"handle"`
		Email  string `json:"email"`
		Status string `json:"status"`
	} `json:"user"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
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

type UserResponse struct {
	ID     string `json:"id"`
	Handle string `json:"handle"`
	Email  string `json:"email"`
}
