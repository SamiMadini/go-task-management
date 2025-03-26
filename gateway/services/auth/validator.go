package auth

import (
	"regexp"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (i *SignUpInput) Validate() []ValidationError {
	var errors []ValidationError

	if i.Handle == "" {
		errors = append(errors, ValidationError{
			Field:   "handle",
			Message: "Handle is required",
		})
	}

	if len(i.Handle) > 50 {
		errors = append(errors, ValidationError{
			Field:   "handle",
			Message: "Handle must be less than 50 characters",
		})
	}

	if i.Email == "" {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "Email is required",
		})
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(i.Email) {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "Invalid email format",
		})
	}

	if i.Password == "" {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password is required",
		})
	}

	if len(i.Password) < 8 {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password must be at least 8 characters",
		})
	}

	return errors
}

func (i *SignInInput) Validate() []ValidationError {
	var errors []ValidationError

	if i.Email == "" {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "Email is required",
		})
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(i.Email) {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "Invalid email format",
		})
	}

	if i.Password == "" {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password is required",
		})
	}

	return errors
}
