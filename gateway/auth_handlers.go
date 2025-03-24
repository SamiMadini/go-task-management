package main

import (
	"context"
	"encoding/base64"
	"log"
	"math/rand"
	"net/http"
	"time"

	commons "sama/go-task-management/commons"
	pb "sama/go-task-management/commons/api"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User         SignupResponse   `json:"user"`
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	ExpiresIn    int64            `json:"expires_in"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// @Summary Sign in user
// @Description Authenticates a user and returns access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param user body LoginRequest true "Sign in credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/signin [post]
func (h *handler) Signin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := commons.ReadJSON(r, &req); err != nil {
		log.Printf("Invalid signin request: %v", err)
		commons.WriteJSONError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.Email == "" || req.Password == "" {
		commons.WriteJSONError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	user, err := h.userRepository.GetByEmail(req.Email)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	if user.ID == "" {
		commons.WriteJSONError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if user.Status != "active" {
		commons.WriteJSONError(w, http.StatusUnauthorized, "Account is not active")
		return
	}

	salt, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		log.Printf("Error decoding salt: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password+string(salt)))
	if err != nil {
		commons.WriteJSONError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	tokens, err := generateTokens(user.ID)
	if err != nil {
		log.Printf("Error generating tokens: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	response := LoginResponse{
		User: SignupResponse{
			ID:     user.ID,
			Handle: user.Handle,
			Email:  user.Email,
			Status: user.Status,
		},
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	}

	commons.WriteJSON(w, http.StatusOK, response)
}

// @Summary Refresh access token
// @Description Generates new access and refresh tokens using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/refresh [post]
func (h *handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := commons.ReadJSON(r, &req); err != nil {
		log.Printf("Invalid refresh token request: %v", err)
		commons.WriteJSONError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.RefreshToken == "" {
		commons.WriteJSONError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	token, err := parseToken(req.RefreshToken)
	if err != nil {
		commons.WriteJSONError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		commons.WriteJSONError(w, http.StatusUnauthorized, "Invalid refresh token claims")
		return
	}

	tokens, err := generateTokens(claims.UserID)
	if err != nil {
		log.Printf("Error generating new tokens: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	commons.WriteJSON(w, http.StatusOK, tokens)
}

// @Summary Sign out user
// @Description Logs out the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} ErrorResponse
// @Router /auth/signout [post]
func (h *handler) Signout(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r)
	if userID == "" {
		commons.WriteJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	log.Printf("User %s signed out successfully", userID)
	commons.WriteJSON(w, http.StatusOK, map[string]string{"message": "Signed out successfully"})
}

// @Summary Request password reset
// @Description Initiates a password reset process by sending a reset link to the user's email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Password reset request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/forgot-password [post]
func (h *handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := commons.ReadJSON(r, &req); err != nil {
		log.Printf("Invalid forgot password request: %v", err)
		commons.WriteJSONError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.Email == "" {
		commons.WriteJSONError(w, http.StatusBadRequest, "Email is required")
		return
	}

	user, err := h.userRepository.GetByEmail(req.Email)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	if user.ID == "" {
		// Return success even if user not found to prevent email enumeration
		commons.WriteJSON(w, http.StatusOK, map[string]string{
			"message": "If the email exists, a password reset link will be sent",
		})
		return
	}

	if user.Status != "active" {
		commons.WriteJSONError(w, http.StatusBadRequest, "Account is not active")
		return
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		log.Printf("Error generating reset token: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)


	now := time.Now()
	resetToken := commons.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: now.Add(24 * time.Hour), // Token expires in 24 hours
		CreatedAt: now,
	}

	_, err = h.passwordResetTokenRepository.Create(resetToken)
	if err != nil {
		log.Printf("Error creating password reset token: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	correlationId := uuid.New().String()
	event := createTaskSystemEvent(
		user.ID,
		correlationId,
		"API Gateway",
		"api:event:password-reset-requested",
		"Password reset requested",
		"{}",
	)

	_, err = h.taskSystemEventRepository.Create(event, 1)
	if err != nil {
		log.Printf("Error creating password reset event: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), notificationTimeout)
	defer cancel()

	_, err = h.notificationServiceClient.SendNotification(
		ctx,
		&pb.SendNotificationRequest{
			TaskId:        user.ID,
			CorrelationId: correlationId,
			Types:         []pb.NotificationType{pb.NotificationType_EMAIL},
		},
		grpc.FailFastCallOption{},
	)

	if err != nil {
		log.Printf("Failed to send password reset notification: %v", err)
	}

	commons.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "If the email exists, a password reset link will be sent",
	})
}

// @Summary Reset password
// @Description Resets the user's password using a valid reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Password reset details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/reset-password [post]
func (h *handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := commons.ReadJSON(r, &req); err != nil {
		log.Printf("Invalid reset password request: %v", err)
		commons.WriteJSONError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.Token == "" || req.NewPassword == "" {
		commons.WriteJSONError(w, http.StatusBadRequest, "Token and new password are required")
		return
	}

	resetToken, err := h.passwordResetTokenRepository.GetByToken(req.Token)
	if err != nil {
		log.Printf("Error retrieving reset token: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	if resetToken.ID == "" {
		commons.WriteJSONError(w, http.StatusBadRequest, "Invalid or expired reset token")
		return
	}

	user, err := h.userRepository.GetByID(resetToken.UserID)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	if user.ID == "" {
		commons.WriteJSONError(w, http.StatusBadRequest, "User not found")
		return
	}

	if user.Status != "active" {
		commons.WriteJSONError(w, http.StatusBadRequest, "Account is not active")
		return
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		log.Printf("Error generating salt: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword+string(salt)), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	_, err = h.userRepository.UpdatePassword(user.ID, string(hashedPassword), base64.StdEncoding.EncodeToString(salt))
	if err != nil {
		log.Printf("Error updating password: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	err = h.passwordResetTokenRepository.MarkAsUsed(resetToken.ID)
	if err != nil {
		log.Printf("Error marking reset token as used: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	commons.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Password has been reset successfully",
	})
}
