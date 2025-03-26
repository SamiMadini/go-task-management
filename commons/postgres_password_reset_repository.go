package commons

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type PasswordResetTokenRepositoryInterface interface {
	Create(token PasswordResetToken) (PasswordResetToken, error)
	GetByToken(token string) (PasswordResetToken, error)
	MarkAsUsed(token string) error
	DeleteExpired() error
}

type PostgresPasswordResetTokenRepository struct {
	DB *sql.DB
}

func NewPostgresPasswordResetTokenRepository(db *sql.DB) *PostgresPasswordResetTokenRepository {
	return &PostgresPasswordResetTokenRepository{DB: db}
}

func (r *PostgresPasswordResetTokenRepository) Create(token PasswordResetToken) (PasswordResetToken, error) {
	dbToken := &DBPasswordResetToken{}
	dbToken.FromPasswordResetToken(token)

	if dbToken.ID == "" {
		dbToken.ID = uuid.New().String()
	}

	if dbToken.CreatedAt.IsZero() {
		dbToken.CreatedAt = time.Now()
	}

	_, err := r.DB.Exec(`
		INSERT INTO password_reset_tokens (id, user_id, token, expires_at, used, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		dbToken.ID,
		dbToken.UserID,
		dbToken.Token,
		dbToken.ExpiresAt,
		dbToken.Used,
		dbToken.CreatedAt,
	)

	if err != nil {
		return PasswordResetToken{}, err
	}

	return dbToken.ToPasswordResetToken(), nil
}

func (r *PostgresPasswordResetTokenRepository) GetByToken(token string) (PasswordResetToken, error) {
	var dbToken DBPasswordResetToken

	err := r.DB.QueryRow(`
		SELECT id, user_id, token, expires_at, used, created_at
		FROM password_reset_tokens
		WHERE token = $1 AND used = false AND expires_at > NOW()
	`, token).Scan(
		&dbToken.ID,
		&dbToken.UserID,
		&dbToken.Token,
		&dbToken.ExpiresAt,
		&dbToken.Used,
		&dbToken.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return PasswordResetToken{}, ErrNotFound
		}
		return PasswordResetToken{}, err
	}

	return dbToken.ToPasswordResetToken(), nil
}

func (r *PostgresPasswordResetTokenRepository) MarkAsUsed(token string) error {
	result, err := r.DB.Exec(`
		UPDATE password_reset_tokens
		SET used = true
		WHERE token = $1 AND used = false AND expires_at > NOW()
	`, token)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	log.Println("Password reset token marked as used successfully")
	return nil
}

func (r *PostgresPasswordResetTokenRepository) DeleteExpired() error {
	_, err := r.DB.Exec(`
		DELETE FROM password_reset_tokens
		WHERE expires_at <= NOW() OR used = true
	`)

	if err != nil {
		return err
	}

	log.Println("Expired password reset tokens deleted successfully")
	return nil
}
