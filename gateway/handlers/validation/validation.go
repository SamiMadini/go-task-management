package validation

import (
	"regexp"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateEmail(email string) *ValidationError {
	if email == "" {
		return &ValidationError{
			Field:   "email",
			Message: "Email is required",
		}
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return &ValidationError{
			Field:   "email",
			Message: "Invalid email format",
		}
	}

	return nil
}

func ValidatePassword(password, field string) *ValidationError {
	if password == "" {
		return &ValidationError{
			Field:   field,
			Message: field + " is required",
		}
	}

	if len(password) < 8 {
		return &ValidationError{
			Field:   field,
			Message: field + " must be at least 8 characters",
		}
	}

	return nil
}

func ValidateHandle(handle string) *ValidationError {
	if handle == "" {
		return &ValidationError{
			Field:   "handle",
			Message: "Handle is required",
		}
	}

	if len(handle) > 50 {
		return &ValidationError{
			Field:   "handle",
			Message: "Handle must be less than 50 characters",
		}
	}

	return nil
}

func ValidateTitle(title string) *ValidationError {
	if title == "" {
		return &ValidationError{
			Field:   "title",
			Message: "Title is required",
		}
	}

	if len(title) > 200 {
		return &ValidationError{
			Field:   "title",
			Message: "Title must be less than 200 characters",
		}
	}

	return nil
}

func ValidateDescription(description string) *ValidationError {
	if len(description) > 1000 {
		return &ValidationError{
			Field:   "description",
			Message: "Description must be less than 1000 characters",
		}
	}

	return nil
}
