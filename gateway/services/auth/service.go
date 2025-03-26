package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"sama/go-task-management/commons"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetByID(id string) (commons.User, error)
	GetByEmail(email string) (commons.User, error)
	Create(user commons.User) (commons.User, error)
	Update(user commons.User) error
	UpdatePassword(userID string, hashedPassword string) error
}

type PasswordResetTokenRepository interface {
	Create(token commons.PasswordResetToken) (commons.PasswordResetToken, error)
	GetByToken(token string) (commons.PasswordResetToken, error)
	MarkAsUsed(token string) error
	DeleteExpired() error
}

type Service struct {
	logger                 commons.Logger
	jwtSecret              string
	userRepo               UserRepository
	passwordResetTokenRepo PasswordResetTokenRepository
}

func NewService(logger commons.Logger, jwtSecret string, userRepo UserRepository, passwordResetTokenRepo PasswordResetTokenRepository) *Service {
	return &Service{
		logger:                 logger,
		jwtSecret:              jwtSecret,
		userRepo:               userRepo,
		passwordResetTokenRepo: passwordResetTokenRepo,
	}
}

func (s *Service) SignUp(ctx context.Context, input SignUpInput) (*AuthResponse, error) {
	_, err := s.userRepo.GetByEmail(input.Email)
	if err == nil {
		return nil, commons.ErrEmailTaken
	}

	salt := generateSalt()
	hashedPassword := hashPassword(input.Password, salt)

	user := commons.User{
		Handle:         input.Handle,
		Email:          input.Email,
		HashedPassword: hashedPassword,
		Salt:          salt,
		Status:        "ACTIVE",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	tokens := s.generateTokens(createdUser.ID)

	return &AuthResponse{
		User: UserResponse{
			ID:     createdUser.ID,
			Handle: createdUser.Handle,
			Email:  createdUser.Email,
			Status: createdUser.Status,
		},
		Token: &tokens,
	}, nil
}

func (s *Service) SignIn(ctx context.Context, input SignInInput) (*AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, commons.ErrInvalidCredentials
	}

	if !verifyPassword(input.Password, user.HashedPassword, user.Salt) {
		return nil, commons.ErrInvalidCredentials
	}

	tokens := s.generateTokens(user.ID)

	return &AuthResponse{
		User: UserResponse{
			ID:     user.ID,
			Handle: user.Handle,
			Email:  user.Email,
			Status: user.Status,
		},
		Token: &tokens,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, userID string) (*TokenResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, commons.ErrInvalidCredentials
	}

	tokens := s.generateTokens(user.ID)
	return &tokens, nil
}

func (s *Service) ForgotPassword(ctx context.Context, input ForgotPasswordInput) error {
	user, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		return commons.ErrNotFound
	}

	token := commons.PasswordResetToken{
		UserID:    user.ID,
		Token:     generateToken(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	_, err = s.passwordResetTokenRepo.Create(token)
	if err != nil {
		return err
	}

	// TODO: Send email with reset link

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	resetToken, err := s.passwordResetTokenRepo.GetByToken(input.Token)
	if err != nil {
		return commons.ErrNotFound
	}

	if resetToken.Used {
		return commons.ErrInvalidInput
	}

	if resetToken.ExpiresAt.Before(time.Now()) {
		return commons.ErrInvalidInput
	}

	salt := generateSalt()
	hashedPassword := hashPassword(input.NewPassword, salt)

	err = s.userRepo.UpdatePassword(resetToken.UserID, hashedPassword)
	if err != nil {
		return err
	}

	err = s.passwordResetTokenRepo.MarkAsUsed(input.Token)
	if err != nil {
		s.logger.Printf("Error marking token as used: %v", err)
	}

	return nil
}

func (s *Service) generateTokens(userID string) TokenResponse {
	now := time.Now()
	accessTokenExpiry := now.Add(time.Hour)
	refreshTokenExpiry := now.Add(7 * 24 * time.Hour)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": accessTokenExpiry.Unix(),
		"iat": now.Unix(),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": refreshTokenExpiry.Unix(),
		"iat": now.Unix(),
	})

	signedAccessToken, _ := accessToken.SignedString([]byte(s.jwtSecret))
	signedRefreshToken, _ := refreshToken.SignedString([]byte(s.jwtSecret))

	return TokenResponse{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
		ExpiresIn:    int64(time.Until(accessTokenExpiry).Seconds()),
	}
}

func generateSalt() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func generateToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func hashPassword(password, salt string) string {
	saltedPassword := password + salt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashedBytes)
}

func verifyPassword(password, hashedPassword, salt string) bool {
	saltedPassword := password + salt
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(saltedPassword))
	return err == nil
}
