package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	commons "sama/go-task-management/commons"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Handle   string `json:"handle"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	ID     string `json:"id"`
	Handle string `json:"handle"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

// @Summary Sign up a new user
// @Description Creates a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body SignupRequest true "User details"
// @Success 201 {object} SignupResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/signup [post]
func (h *handler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		commons.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req SignupRequest
	if err := commons.ReadJSON(r, &req); err != nil {
		log.Printf("Invalid signup request: %v", err)
		commons.WriteJSONError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.Handle == "" || req.Email == "" || req.Password == "" {
		commons.WriteJSONError(w, http.StatusBadRequest, "All fields are required")
		return
	}

	existingUser, err := h.userRepository.GetByEmail(req.Email)
	if err != nil {
		log.Printf("Error checking existing user: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}
	if existingUser.ID != "" {
		commons.WriteJSONError(w, http.StatusConflict, "Email already registered")
		return
	}

	existingUser, err = h.userRepository.GetByHandle(req.Handle)
	if err != nil {
		log.Printf("Error checking existing user: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}
	if existingUser.ID != "" {
		commons.WriteJSONError(w, http.StatusConflict, "Handle already taken")
		return
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		log.Printf("Error generating salt: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password+string(salt)), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	now := time.Now()
	user := commons.User{
		ID:           uuid.New().String(),
		Handle:       req.Handle,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Salt:         base64.StdEncoding.EncodeToString(salt),
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result, err := h.userRepository.Create(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	response := SignupResponse{
		ID:     result.ID,
		Handle: result.Handle,
		Email:  result.Email,
		Status: result.Status,
	}

	commons.WriteJSON(w, http.StatusCreated, response)
}
