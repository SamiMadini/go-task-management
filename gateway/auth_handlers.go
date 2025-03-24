package main

import (
	"encoding/base64"
	"log"
	"net/http"

	commons "sama/go-task-management/commons"

	"golang.org/x/crypto/bcrypt"
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
