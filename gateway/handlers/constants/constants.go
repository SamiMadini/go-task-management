package constants

import "time"

const (
	TaskStatusTodo       = "TODO"
	TaskStatusInProgress = "IN_PROGRESS"
	TaskStatusDone       = "DONE"
)

const (
	TaskPriorityLow    = 1
	TaskPriorityMedium = 2
	TaskPriorityHigh   = 3
)

const (
	ErrCodeValidation   = "VALIDATION_ERROR"
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeForbidden    = "FORBIDDEN"
	ErrCodeInternal     = "INTERNAL_ERROR"
	ErrCodeBadRequest   = "BAD_REQUEST"
)

const (
	NotificationTimeout = 5 * time.Second
)
