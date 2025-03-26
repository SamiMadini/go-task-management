package types

// ErrorResponse represents a standardized error response
// @Description Error response model with detailed information about the error
// @Description Common error codes:
// @Description - BAD_REQUEST: Invalid input or request format
// @Description - UNAUTHORIZED: Authentication required
// @Description - FORBIDDEN: Permission denied
// @Description - NOT_FOUND: Resource not found
// @Description - INTERNAL_ERROR: Server error
// @Description - VALIDATION_ERROR: Input validation failed
type ErrorResponse struct {
	Code             string            `json:"code" example:"BAD_REQUEST" enums:"BAD_REQUEST,UNAUTHORIZED,FORBIDDEN,NOT_FOUND,INTERNAL_ERROR,VALIDATION_ERROR"`
	Message          string            `json:"message" example:"Invalid task format"`
	Details          string            `json:"details,omitempty" example:"Task ID must be a valid UUID"`
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
}

// ValidationError represents a field-level validation error
// @Description Validation error for a specific field
type ValidationError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"Email address is invalid"`
}
