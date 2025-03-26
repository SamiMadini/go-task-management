package commons

import (
	"net/http"
)

var (
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
	}
	InternalServerErrorHandler = func(w http.ResponseWriter) {
		WriteJSONError(w, http.StatusInternalServerError, "An unexpected error occurred")
	}
	NotFoundErrorHandler = func(w http.ResponseWriter, err error) {
		WriteJSONError(w, http.StatusNotFound, err.Error())
	}

	ErrForbidden = NewError("FORBIDDEN", "Access forbidden")

	ErrNotFound = NewError("NOT_FOUND", "Resource not found")

	ErrInvalidInput = NewError("INVALID_INPUT", "Invalid input")

	ErrUnauthorized = NewError("UNAUTHORIZED", "Unauthorized")

	ErrInternal = NewError("INTERNAL_ERROR", "Internal error")

	ErrEmailTaken = NewError("EMAIL_TAKEN", "Email already taken")

	ErrInvalidCredentials = NewError("INVALID_CREDENTIALS", "Invalid credentials")
)
