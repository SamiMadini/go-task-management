package task

import (
	"time"
)

type CreateTaskInput struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    int        `json:"priority"`
	DueDate     time.Time  `json:"due_date"`
	CreatorID   string     `json:"creator_id"`
	AssigneeID  *string    `json:"assignee_id,omitempty"`
}

type UpdateTaskInput struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    int        `json:"priority"`
	DueDate     time.Time  `json:"due_date"`
	AssigneeID  *string    `json:"assignee_id,omitempty"`
	UserID      string     `json:"user_id"` // The ID of the user making the update
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (i *CreateTaskInput) Validate() []ValidationError {
	var errors []ValidationError

	if i.Title == "" {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "Title is required",
		})
	}

	if len(i.Title) > 200 {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "Title must be less than 200 characters",
		})
	}

	if len(i.Description) > 1000 {
		errors = append(errors, ValidationError{
			Field:   "description",
			Message: "Description must be less than 1000 characters",
		})
	}

	if i.Status != "TODO" && i.Status != "IN_PROGRESS" && i.Status != "DONE" {
		errors = append(errors, ValidationError{
			Field:   "status",
			Message: "Invalid status. Must be one of: TODO, IN_PROGRESS, DONE",
		})
	}

	if i.Priority < 1 || i.Priority > 3 {
		errors = append(errors, ValidationError{
			Field:   "priority",
			Message: "Priority must be between 1 and 3",
		})
	}

	if i.DueDate.Before(time.Now()) {
		errors = append(errors, ValidationError{
			Field:   "due_date",
			Message: "Due date must be in the future",
		})
	}

	return errors
}

func (i *UpdateTaskInput) Validate() []ValidationError {
	var errors []ValidationError

	if i.Title != "" && len(i.Title) > 200 {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "Title must be less than 200 characters",
		})
	}

	if i.Description != "" && len(i.Description) > 1000 {
		errors = append(errors, ValidationError{
			Field:   "description",
			Message: "Description must be less than 1000 characters",
		})
	}

	if i.Status != "" && i.Status != "TODO" && i.Status != "IN_PROGRESS" && i.Status != "DONE" {
		errors = append(errors, ValidationError{
			Field:   "status",
			Message: "Invalid status. Must be one of: TODO, IN_PROGRESS, DONE",
		})
	}

	if i.Priority != 0 && (i.Priority < 1 || i.Priority > 3) {
		errors = append(errors, ValidationError{
			Field:   "priority",
			Message: "Priority must be between 1 and 3",
		})
	}

	if !i.DueDate.IsZero() && i.DueDate.Before(time.Now()) {
		errors = append(errors, ValidationError{
			Field:   "due_date",
			Message: "Due date must be in the future",
		})
	}

	return errors
}
