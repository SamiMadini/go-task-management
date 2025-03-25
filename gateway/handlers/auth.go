package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"sama/go-task-management/commons"
	"sama/go-task-management/gateway/middleware"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Handle   string `json:"handle"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  SignupResponse `json:"user"`
	Token TokenResponse  `json:"token"`
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

type AuthHandler struct {
	*BaseHandler
}

func NewAuthHandler(base *BaseHandler) *AuthHandler {
	return &AuthHandler{BaseHandler: base}
}

func (h *AuthHandler) generateTokens(userID string) TokenResponse {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &middleware.AuthClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &middleware.AuthClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	})

	accessTokenString, _ := accessToken.SignedString([]byte(h.jwtSecret))
	refreshTokenString, _ := refreshToken.SignedString([]byte(h.jwtSecret))

	return TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    15 * 60, // 15 minutes in seconds
	}
}

func (h *AuthHandler) hashPassword(password, salt string) string {
	saltedPassword := password + salt

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return ""
	}

	return string(hashedBytes)
}

func (h *AuthHandler) comparePasswords(hashedPassword, password, salt string) bool {
	saltedPassword := password + salt
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(saltedPassword))
	if err != nil {
		log.Printf("Password comparison failed: %v", err)
		return false
	}
	log.Printf("Password comparison successful")
	return true
}

func (h *AuthHandler) generateSalt() string {
	// Generate a shorter salt (16 bytes = 32 characters when base64 encoded)
	saltBytes := make([]byte, 16)
	_, err := rand.Read(saltBytes)
	if err != nil {
		log.Printf("Error generating salt: %v", err)
		return ""
	}

	return base64.URLEncoding.EncodeToString(saltBytes)
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
func (h *AuthHandler) Signin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := h.decodeJSON(r, &req); err != nil {
		log.Printf("Error decoding request: %v", err)
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Printf("Attempting to sign in user with email: %s", req.Email)

	user, err := h.userRepository.GetByEmail(req.Email)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		h.respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if user.ID == "" {
		log.Printf("User not found with email: %s", req.Email)
		h.respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	log.Printf("Found user: ID=%s, Email=%s", user.ID, user.Email)

	if !h.comparePasswords(user.PasswordHash, req.Password, user.Salt) {
		log.Printf("Invalid password for user: %s", user.Email)
		h.respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	log.Printf("Password verified for user: %s", user.Email)

	tokenResponse := h.generateTokens(user.ID)

	userResponse := SignupResponse{}
	userResponse.User.ID = user.ID
	userResponse.User.Handle = user.Handle
	userResponse.User.Email = user.Email
	userResponse.User.Status = user.Status

	response := LoginResponse{
		User:  userResponse,
		Token: tokenResponse,
	}

	h.respondWithJSON(w, http.StatusOK, response)
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
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	token, err := middleware.ParseToken(req.RefreshToken, h.jwtSecret)
	if err != nil {
		h.respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	claims, ok := token.Claims.(*middleware.AuthClaims)
	if !ok || !token.Valid {
		h.respondWithError(w, http.StatusUnauthorized, "Invalid refresh token claims")
		return
	}

	tokenResponse := h.generateTokens(claims.UserID)
	h.respondWithJSON(w, http.StatusOK, tokenResponse)
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
func (h *AuthHandler) Signout(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// TODO: Invalidate the refresh token here or add it to a blacklist

	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully signed out"})
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
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.userRepository.GetByEmail(req.Email)
	if err != nil {
		// Avoid enumeration attacks
		h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "If the email exists, a password reset link has been sent"})
		return
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	resetToken := commons.PasswordResetToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
		Used:      false,
		CreatedAt: time.Now(),
	}

	_, err = h.passwordResetTokenRepository.Create(resetToken)
	if err != nil {
		log.Printf("Error creating password reset token: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to process request")
		return
	}

	// TODO: Send email with reset link
	resetLink := "http://localhost:3000/reset-password?token=" + token
	h.respondWithJSON(w, http.StatusOK, map[string]string{"reset_link": resetLink})
}

// @Summary Reset password
// @Description Resets the user's password using a valid reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Password reset request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resetToken, err := h.passwordResetTokenRepository.GetByToken(req.Token)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid or expired reset token")
		return
	}

	if resetToken.Used || time.Now().After(resetToken.ExpiresAt) {
		h.respondWithError(w, http.StatusBadRequest, "Invalid or expired reset token")
		return
	}

	salt := h.generateSalt()
	if salt == "" {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to generate salt")
		return
	}

	hashedPassword := h.hashPassword(req.NewPassword, salt)
	if hashedPassword == "" {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	_, err = h.userRepository.UpdatePassword(resetToken.UserID, hashedPassword, salt)
	if err != nil {
		log.Printf("Error updating password: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to reset password")
		return
	}

	err = h.passwordResetTokenRepository.MarkAsUsed(resetToken.ID)
	if err != nil {
		log.Printf("Error marking reset token as used: %v", err)
		// Don't return an error to the user since the password was successfully reset
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "Password successfully reset"})
}

// @Summary Sign up a new user
// @Description Creates a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body SignupRequest true "User registration data"
// @Success 201 {object} SignupResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/signup [post]
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req SignupRequest
	if err := h.decodeJSON(r, &req); err != nil {
		log.Printf("Invalid signup request: %v", err)
		h.respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.Handle == "" || req.Email == "" || req.Password == "" {
		h.respondWithError(w, http.StatusBadRequest, "All fields are required")
		return
	}

	existingUser, err := h.userRepository.GetByEmail(req.Email)
	if err != nil {
		log.Printf("Error checking existing user: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if existingUser.ID != "" {
		h.respondWithError(w, http.StatusConflict, "Email already registered")
		return
	}

	existingUser, err = h.userRepository.GetByHandle(req.Handle)
	if err != nil {
		log.Printf("Error checking existing user: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if existingUser.ID != "" {
		h.respondWithError(w, http.StatusConflict, "Handle already taken")
		return
	}

	salt := h.generateSalt()
	if salt == "" {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to generate salt")
		return
	}

	hashedPassword := h.hashPassword(req.Password, salt)
	if hashedPassword == "" {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	now := time.Now()
	user := commons.User{
		ID:           uuid.New().String(),
		Handle:       req.Handle,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Salt:         salt,
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result, err := h.userRepository.Create(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response := SignupResponse{
		User: struct {
			ID     string `json:"id"`
			Handle string `json:"handle"`
			Email  string `json:"email"`
			Status string `json:"status"`
		}{
			ID:     result.ID,
			Handle: result.Handle,
			Email:  result.Email,
			Status: result.Status,
		},
	}

	h.respondWithJSON(w, http.StatusCreated, response)
}
