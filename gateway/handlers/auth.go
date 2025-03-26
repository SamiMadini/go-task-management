package handlers

import (
	"encoding/json"
	"net/http"

	"sama/go-task-management/commons"
	"sama/go-task-management/gateway/handlers/constants"
	"sama/go-task-management/gateway/handlers/validation"
	"sama/go-task-management/gateway/services/auth"

	"github.com/golang-jwt/jwt/v5"
)

type SignUpRequest struct {
	Handle   string `json:"handle"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *SignUpRequest) Validate() []validation.ValidationError {
	var errors []validation.ValidationError

	if handleErr := validation.ValidateHandle(r.Handle); handleErr != nil {
		errors = append(errors, *handleErr)
	}

	if emailErr := validation.ValidateEmail(r.Email); emailErr != nil {
		errors = append(errors, *emailErr)
	}

	if passwordErr := validation.ValidatePassword(r.Password, "password"); passwordErr != nil {
		errors = append(errors, *passwordErr)
	}

	return errors
}

type SignUpResponse struct {
	User  auth.UserResponse `json:"user"`
	Token TokenResponse     `json:"token"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *SignInRequest) Validate() []validation.ValidationError {
	var errors []validation.ValidationError

	if emailErr := validation.ValidateEmail(r.Email); emailErr != nil {
		errors = append(errors, *emailErr)
	}

	if passwordErr := validation.ValidatePassword(r.Password, "password"); passwordErr != nil {
		errors = append(errors, *passwordErr)
	}

	return errors
}

type SignInResponse struct {
	User  auth.UserResponse `json:"user"`
	Token TokenResponse     `json:"token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r *RefreshTokenRequest) Validate() []validation.ValidationError {
	var errors []validation.ValidationError

	if r.RefreshToken == "" {
		errors = append(errors, validation.ValidationError{
			Field:   "refresh_token",
			Message: "Refresh token is required",
		})
	}

	return errors
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func (r *ForgotPasswordRequest) Validate() []validation.ValidationError {
	var errors []validation.ValidationError

	if emailErr := validation.ValidateEmail(r.Email); emailErr != nil {
		errors = append(errors, *emailErr)
	}

	return errors
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func (r *ResetPasswordRequest) Validate() []validation.ValidationError {
	var errors []validation.ValidationError

	if r.Token == "" {
		errors = append(errors, validation.ValidationError{
			Field:   "token",
			Message: "Token is required",
		})
	}

	if passwordErr := validation.ValidatePassword(r.NewPassword, "new_password"); passwordErr != nil {
		errors = append(errors, *passwordErr)
	}

	return errors
}

type AuthHandler struct {
	*BaseHandler
	authService *auth.Service
}

func NewAuthHandler(base *BaseHandler, authService *auth.Service) *AuthHandler {
	return &AuthHandler{
		BaseHandler: base,
		authService: authService,
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var input SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondWithError(w, http.StatusBadRequest, constants.ErrCodeBadRequest, "Invalid request payload", err.Error())
		return
	}

	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		h.respondWithValidationErrors(w, validationErrors)
		return
	}

	response, err := h.authService.SignUp(r.Context(), auth.SignUpInput{
		Handle:   input.Handle,
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		switch err {
		case commons.ErrInvalidInput:
			h.respondWithError(w, http.StatusBadRequest, constants.ErrCodeBadRequest, "Invalid input", err.Error())
		case commons.ErrEmailTaken:
			h.respondWithError(w, http.StatusConflict, constants.ErrCodeBadRequest, "Email already taken", err.Error())
		default:
			h.respondWithError(w, http.StatusInternalServerError, constants.ErrCodeInternal, "Failed to create user", err.Error())
		}
		return
	}

	h.respondWithJSON(w, http.StatusCreated, StandardResponse{
		Success: true,
		Data:    response,
	})
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var input SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondWithError(w, http.StatusBadRequest, constants.ErrCodeBadRequest, "Invalid request payload", err.Error())
		return
	}

	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		h.respondWithValidationErrors(w, validationErrors)
		return
	}

	response, err := h.authService.SignIn(r.Context(), auth.SignInInput{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		switch err {
		case commons.ErrInvalidCredentials:
			h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Invalid credentials", "")
		default:
			h.respondWithError(w, http.StatusInternalServerError, constants.ErrCodeInternal, "Failed to sign in", err.Error())
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, StandardResponse{
		Success: true,
		Data:    response,
	})
}

func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	h.respondWithJSON(w, http.StatusOK, StandardResponse{
		Success: true,
		Data: map[string]string{
			"message": "Signed out successfully",
		},
	})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var input RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondWithError(w, http.StatusBadRequest, constants.ErrCodeBadRequest, "Invalid request payload", err.Error())
		return
	}

	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		h.respondWithValidationErrors(w, validationErrors)
		return
	}

	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Invalid refresh token", "")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Invalid refresh token", "")
		return
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Invalid refresh token", "")
		return
	}

	response, err := h.authService.RefreshToken(r.Context(), userID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, constants.ErrCodeInternal, "Failed to refresh token", err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, StandardResponse{
		Success: true,
		Data:    response,
	})
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var input ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondWithError(w, http.StatusBadRequest, constants.ErrCodeBadRequest, "Invalid request payload", err.Error())
		return
	}

	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		h.respondWithValidationErrors(w, validationErrors)
		return
	}

	err := h.authService.ForgotPassword(r.Context(), auth.ForgotPasswordInput{
		Email: input.Email,
	})
	if err != nil {
		switch err {
		case commons.ErrNotFound:
			h.respondWithError(w, http.StatusNotFound, constants.ErrCodeNotFound, "User not found", "")
		default:
			h.respondWithError(w, http.StatusInternalServerError, constants.ErrCodeInternal, "Failed to process forgot password request", err.Error())
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, StandardResponse{
		Success: true,
		Data: map[string]string{
			"message": "Password reset instructions sent to your email",
		},
	})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var input ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondWithError(w, http.StatusBadRequest, constants.ErrCodeBadRequest, "Invalid request payload", err.Error())
		return
	}

	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		h.respondWithValidationErrors(w, validationErrors)
		return
	}

	err := h.authService.ResetPassword(r.Context(), auth.ResetPasswordInput{
		Token:       input.Token,
		NewPassword: input.NewPassword,
	})
	if err != nil {
		switch err {
		case commons.ErrInvalidInput:
			h.respondWithError(w, http.StatusBadRequest, constants.ErrCodeBadRequest, "Invalid or expired token", "")
		default:
			h.respondWithError(w, http.StatusInternalServerError, constants.ErrCodeInternal, "Failed to reset password", err.Error())
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, StandardResponse{
		Success: true,
		Data: map[string]string{
			"message": "Password reset successfully",
		},
	})
}
