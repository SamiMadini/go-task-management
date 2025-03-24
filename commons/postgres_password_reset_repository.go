package commons

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

// @Description PasswordResetToken model
type PasswordResetToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

type PasswordResetTokenRepositoryInterface interface {
	Create(token PasswordResetToken) (PasswordResetToken, error)
	GetByToken(token string) (PasswordResetToken, error)
	MarkAsUsed(id string) error
	DeleteExpired() error
}

type PostgresPasswordResetTokenRepository struct {
	DB *sql.DB
}

func NewPostgresPasswordResetTokenRepository(db *sql.DB) *PostgresPasswordResetTokenRepository {
	return &PostgresPasswordResetTokenRepository{DB: db}
}

func (r *PostgresPasswordResetTokenRepository) Create(token PasswordResetToken) (PasswordResetToken, error) {
	if token.ID == "" {
		token.ID = uuid.New().String()
	}

	_, err := r.DB.Exec(`
		INSERT INTO password_reset_tokens (id, user_id, token, expires_at, used, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		token.ID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.Used,
		token.CreatedAt,
	)
	if err != nil {
		log.Printf("Error creating password reset token: %v", err)
		return PasswordResetToken{}, err
	}

	return token, nil
}

func (r *PostgresPasswordResetTokenRepository) GetByToken(token string) (PasswordResetToken, error) {
	var resetToken PasswordResetToken
	err := r.DB.QueryRow(`
		SELECT id, user_id, token, expires_at, used, created_at
		FROM password_reset_tokens
		WHERE token = $1 AND used = false AND expires_at > NOW()
	`, token).Scan(
		&resetToken.ID,
		&resetToken.UserID,
		&resetToken.Token,
		&resetToken.ExpiresAt,
		&resetToken.Used,
		&resetToken.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return PasswordResetToken{}, nil
	}
	if err != nil {
		log.Printf("Error getting password reset token: %v", err)
		return PasswordResetToken{}, err
	}

	return resetToken, nil
}

func (r *PostgresPasswordResetTokenRepository) MarkAsUsed(id string) error {
	_, err := r.DB.Exec(`
		UPDATE password_reset_tokens
		SET used = true
		WHERE id = $1
	`, id)
	if err != nil {
		log.Printf("Error marking password reset token as used: %v", err)
		return err
	}

	return nil
}

func (r *PostgresPasswordResetTokenRepository) DeleteExpired() error {
	_, err := r.DB.Exec(`
		DELETE FROM password_reset_tokens
		WHERE expires_at < NOW() OR used = true
	`)
	if err != nil {
		log.Printf("Error deleting expired password reset tokens: %v", err)
		return err
	}

	return nil
}
