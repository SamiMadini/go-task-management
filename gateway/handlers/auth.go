package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"regexp"
	"time"

	"sama/go-task-management/commons"
	"sama/go-task-management/gateway/middleware"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserResponse struct {
	ID     string `json:"id"`
	Handle string `json:"handle"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type SignUpRequest struct {
	Handle   string `json:"handle"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *SignUpRequest) Validate() []ValidationError {
	var errors []ValidationError

	if r.Handle == "" {
		errors = append(errors, ValidationError{
			Field:   "handle",
			Message: "Handle is required",
		})
	}

	if len(r.Handle) > 50 {
		errors = append(errors, ValidationError{
			Field:   "handle",
			Message: "Handle must be less than 50 characters",
		})
	}

	if r.Email == "" {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "Email is required",
		})
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "Invalid email format",
		})
	}

	if r.Password == "" {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password is required",
		})
	}

	if len(r.Password) < 8 {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password must be at least 8 characters",
		})
	}

	return errors
}

type SignUpResponse struct {
	User  UserResponse   `json:"user"`
	Token TokenResponse  `json:"token"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *SignInRequest) Validate() []ValidationError {
	var errors []ValidationError

	if r.Email == "" {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "Email is required",
		})
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "Invalid email format",
		})
	}

	if r.Password == "" {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password is required",
		})
	}

	return errors
}

type SignInResponse struct {
	User  UserResponse   `json:"user"`
	Token TokenResponse  `json:"token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r *RefreshTokenRequest) Validate() []ValidationError {
	var errors []ValidationError

	if r.RefreshToken == "" {
		errors = append(errors, ValidationError{
			Field:   "refresh_token",
			Message: "Refresh token is required",
		})
	}

	return errors
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func (r *ForgotPasswordRequest) Validate() []ValidationError {
	var errors []ValidationError

	if r.Email == "" {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "Email is required",
		})
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "Invalid email format",
		})
	}

	return errors
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func (r *ResetPasswordRequest) Validate() []ValidationError {
	var errors []ValidationError

	if r.Token == "" {
		errors = append(errors, ValidationError{
			Field:   "token",
			Message: "Token is required",
		})
	}

	if r.NewPassword == "" {
		errors = append(errors, ValidationError{
			Field:   "new_password",
			Message: "New password is required",
		})
	}

	if len(r.NewPassword) < 8 {
		errors = append(errors, ValidationError{
			Field:   "new_password",
			Message: "New password must be at least 8 characters",
		})
	}

	return errors
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
// @Description
// @Description Error scenarios:
// @Description - Invalid credentials: Returns 401 with UNAUTHORIZED code
// @Description - Missing required fields: Returns 400 with VALIDATION_ERROR code
// @Description - Invalid email format: Returns 400 with VALIDATION_ERROR code
// @Tags auth
// @Accept json
// @Produce json
// @Param user body SignInRequest true "Sign in credentials"
// @Success 200 {object} SignInResponse "Successfully authenticated"
// @Failure 400 {object} ErrorResponse "Invalid request format or validation errors"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/signin [post]
func (h *AuthHandler) Signin(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body", err.Error())
		return
	}

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		h.respondWithValidationErrors(w, validationErrors)
		return
	}

	user, err := h.userRepository.GetByEmail(req.Email)
	if err != nil {
		h.logger.Printf("Error getting user: %v", err)
		h.respondWithError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Invalid credentials", "Email or password is incorrect")
		return
	}

	if user.ID == "" {
		h.logger.Printf("User not found with email: %s", req.Email)
		h.respondWithError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Invalid credentials", "Email or password is incorrect")
		return
	}

	if !h.comparePasswords(user.PasswordHash, req.Password, user.Salt) {
		h.logger.Printf("Invalid password for user: %s", user.Email)
		h.respondWithError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Invalid credentials", "Email or password is incorrect")
		return
	}

	tokenResponse := h.generateTokens(user.ID)
	userResponse := UserResponse{ID: user.ID, Handle: user.Handle, Email: user.Email, Status: user.Status}

	h.respondWithJSON(w, http.StatusOK, SignInResponse{
		User:  userResponse,
		Token: tokenResponse,
	})
}

// @Summary Refresh access token
// @Description Generates new access and refresh tokens using a valid refresh token
// @Description
// @Description Error scenarios:
// @Description - Invalid refresh token: Returns 401 with UNAUTHORIZED code
// @Description - Token expired: Returns 401 with UNAUTHORIZED code
// @Description - Token blacklisted: Returns 401 with UNAUTHORIZED code
// @Tags auth
// @Accept json
// @Produce json
// @Param token body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} TokenResponse "New tokens generated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 401 {object} ErrorResponse "Invalid or expired refresh token"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body", err.Error())
		return
	}

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		h.respondWithValidationErrors(w, validationErrors)
		return
	}

	token, err := middleware.ParseToken(req.RefreshToken, h.jwtSecret)
	if err != nil {
		h.respondWithError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Invalid refresh token", err.Error())
		return
	}

	claims, ok := token.Claims.(*middleware.AuthClaims)
	if !ok || !token.Valid {
		h.respondWithError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Invalid refresh token claims", "Token claims are invalid or expired")
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
		h.respondWithError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	// TODO: Implement token blacklisting or revocation

	h.respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Successfully signed out",
	})
}

// @Summary Request password reset
// @Description Initiates a password reset process by sending a reset link to the user's email
// @Description
// @Description Error scenarios:
// @Description - User not found: Returns 404 with NOT_FOUND code
// @Description - Invalid email format: Returns 400 with VALIDATION_ERROR code
// @Description - Email service error: Returns 500 with INTERNAL_ERROR code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Password reset request"
// @Success 200 {object} map[string]string "Reset email sent successfully"
// @Failure 400 {object} ErrorResponse "Invalid email format"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Failed to send reset email"
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body", err.Error())
		return
	}

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		h.respondWithValidationErrors(w, validationErrors)
		return
	}

	user, err := h.userRepository.GetByEmail(req.Email)
	if err != nil {
		h.logger.Printf("Error getting user: %v", err)
		h.respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "If your email is registered, you will receive a password reset link",
		})
		return
	}

	if user.ID == "" {
		h.respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "If your email is registered, you will receive a password reset link",
		})
		return
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour)

	resetToken := commons.PasswordResetToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	_, err = h.passwordResetTokenRepository.Create(resetToken)
	if err != nil {
		h.logger.Printf("Error creating password reset token: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to process request", err.Error())
		return
	}

	// TODO: Send email with reset link

	h.respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "If your email is registered, you will receive a password reset link",
	})
}

// @Summary Reset password
// @Description Resets the user's password using a valid reset token
// @Description
// @Description Error scenarios:
// @Description - Invalid token: Returns 400 with BAD_REQUEST code
// @Description - Token expired: Returns 400 with BAD_REQUEST code
// @Description - Weak password: Returns 400 with VALIDATION_ERROR code
// @Description - Token already used: Returns 400 with BAD_REQUEST code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Password reset request"
// @Success 200 {object} map[string]string "Password reset successful"
// @Failure 400 {object} ErrorResponse "Invalid token or password validation errors"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body", err.Error())
		return
	}

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		h.respondWithValidationErrors(w, validationErrors)
		return
	}

	resetToken, err := h.passwordResetTokenRepository.GetByToken(req.Token)
	if err != nil {
		h.logger.Printf("Error getting reset token: %v", err)
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid or expired token", "The reset token is invalid or has expired")
		return
	}

	if resetToken.ExpiresAt.Before(time.Now()) {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Token expired", "The reset token has expired")
		return
	}

	user, err := h.userRepository.GetByID(resetToken.UserID)
	if err != nil {
		h.logger.Printf("Error getting user: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to process request", err.Error())
		return
	}

	salt := h.generateSalt()
	if salt == "" {
		h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to process request", "Error generating salt")
		return
	}

	passwordHash := h.hashPassword(req.NewPassword, salt)
	if passwordHash == "" {
		h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to process request", "Error hashing password")
		return
	}

	_, err = h.userRepository.UpdatePassword(user.ID, passwordHash, salt)
	if err != nil {
		h.logger.Printf("Error updating user: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to update password", err.Error())
		return
	}

	if err = h.passwordResetTokenRepository.MarkAsUsed(resetToken.ID); err != nil {
		h.logger.Printf("Error marking reset token as used: %v", err)
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Password has been successfully reset",
	})
}

// @Summary Sign up a new user
// @Description Creates a new user account
// @Description
// @Description Error scenarios:
// @Description - Email already registered: Returns 400 with BAD_REQUEST code
// @Description - Handle already taken: Returns 400 with BAD_REQUEST code
// @Description - Invalid email format: Returns 400 with VALIDATION_ERROR code
// @Description - Password too weak: Returns 400 with VALIDATION_ERROR code
// @Description - Missing required fields: Returns 400 with VALIDATION_ERROR code
// @Tags auth
// @Accept json
// @Produce json
// @Param user body SignUpRequest true "User registration data"
// @Success 201 {object} SignUpResponse "Account created successfully"
// @Failure 400 {object} ErrorResponse "Validation errors or duplicate email/handle"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/signup [post]
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body", err.Error())
		return
	}

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		h.respondWithValidationErrors(w, validationErrors)
		return
	}

	existingUser, err := h.userRepository.GetByEmail(req.Email)
	if err == nil && existingUser.ID != "" {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Email already registered", "Please use a different email address")
		return
	}

	existingUser, err = h.userRepository.GetByHandle(req.Handle)
	if err == nil && existingUser.ID != "" {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Handle already taken", "Please choose a different handle")
		return
	}

	salt := h.generateSalt()
	if salt == "" {
		h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to process request", "Error generating salt")
		return
	}

	passwordHash := h.hashPassword(req.Password, salt)
	if passwordHash == "" {
		h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to process request", "Error hashing password")
		return
	}

	now := time.Now()
	user := commons.User{
		ID:           uuid.New().String(),
		Handle:       req.Handle,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Salt:         salt,
		Status:       "ACTIVE",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	createdUser, err := h.userRepository.Create(user)
	if err != nil {
		h.logger.Printf("Error creating user: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to create user", err.Error())
		return
	}

	tokenResponse := h.generateTokens(createdUser.ID)
	userResponse := UserResponse{ID: createdUser.ID, Handle: createdUser.Handle, Email: createdUser.Email, Status: createdUser.Status}

	h.respondWithJSON(w, http.StatusCreated, SignUpResponse{
		User:  userResponse,
		Token: tokenResponse,
	})
}
