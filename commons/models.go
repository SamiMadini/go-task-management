package commons

import (
	"time"
)

type User struct {
	ID             string    `json:"id"`
	Handle         string    `json:"handle"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	Salt           string    `json:"-"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Task struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      string            `json:"status"`
	Priority    int              `json:"priority"`
	DueDate     time.Time         `json:"due_date"`
	CreatorID   string            `json:"creator_id"`
	AssigneeID  *string           `json:"assignee_id,omitempty"`
	EmailSent   bool              `json:"email_sent"`
	InAppSent   bool              `json:"in_app_sent"`
	Deleted     bool              `json:"deleted"`
	DeletedAt   *time.Time        `json:"deleted_at,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Events      []TaskSystemEvent `json:"events,omitempty"`
}

type TaskSystemEvent struct {
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

type Notification struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	IsRead    bool                   `json:"is_read"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

type InAppNotification struct {
	ID          string      `json:"id"`
	UserID      string      `json:"user_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	IsRead      bool        `json:"is_read"`
	ReadAt      *time.Time  `json:"read_at,omitempty"`
	Deleted     bool        `json:"deleted"`
	DeletedAt   *time.Time  `json:"deleted_at,omitempty"`
	UpdatedAt   time.Time   `json:"updated_at"`
	CreatedAt   time.Time   `json:"created_at"`
}

type PasswordResetToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

type GRPCEvent struct {
	TaskId        string
	CorrelationId string
	Types         []string
}

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
