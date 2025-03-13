package commons

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

// @Description InAppNotification model
type InAppNotification struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsRead      bool      `json:"is_read"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type InAppNotificationRepositoryInterface interface {
	GetAll() ([]InAppNotification, error)
	GetByID(id string) (InAppNotification, error)
	Create(inAppNotification InAppNotification) (InAppNotification, error)
	UpdateReadAt(id string) error
	Update(inAppNotification InAppNotification) error
	Delete(id string) error
}

type PostgresInAppNotificationRepository struct {
	DB *sql.DB
}

func NewPostgresInAppNotificationRepository(db *sql.DB) *PostgresInAppNotificationRepository {
	return &PostgresInAppNotificationRepository{DB: db}
}

func (r *PostgresInAppNotificationRepository) GetAll() ([]InAppNotification, error) {
	rows, err := r.DB.Query(`
		SELECT id, title, description, is_read, read_at, created_at, updated_at
		FROM in_app_notifications
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inAppNotifications []InAppNotification
	for rows.Next() {
		var inAppNotification InAppNotification
		var readAt sql.NullTime
		
		err := rows.Scan(
			&inAppNotification.ID,
			&inAppNotification.Title,
			&inAppNotification.Description,
			&inAppNotification.IsRead,
			&readAt,
			&inAppNotification.CreatedAt,
			&inAppNotification.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if readAt.Valid {
			inAppNotification.ReadAt = &readAt.Time
		}
		
		inAppNotifications = append(inAppNotifications, inAppNotification)
	}

	return inAppNotifications, nil
}

func (r *PostgresInAppNotificationRepository) GetByID(id string) (InAppNotification, error) {
	var inAppNotification InAppNotification
	var readAt sql.NullTime
	
	err := r.DB.QueryRow(`
		SELECT id, title, description, is_read, read_at, created_at, updated_at
		FROM in_app_notifications WHERE id = $1
	`, id).Scan(
		&inAppNotification.ID,
		&inAppNotification.Title,
		&inAppNotification.Description,
		&inAppNotification.IsRead,
		&readAt,
		&inAppNotification.CreatedAt,
		&inAppNotification.UpdatedAt,
	)
	
	if err != nil {
		return InAppNotification{}, err
	}

	if readAt.Valid {
		inAppNotification.ReadAt = &readAt.Time
	}
	
	return inAppNotification, nil
}

func (r *PostgresInAppNotificationRepository) Create(inAppNotification InAppNotification) (InAppNotification, error) {
	if inAppNotification.ID == "" {
		inAppNotification.ID = uuid.New().String()
	}

	now := time.Now()
	inAppNotification.CreatedAt = now
	inAppNotification.UpdatedAt = now
	
	_, err := r.DB.Exec(`
		INSERT INTO in_app_notifications (id, title, description, is_read, read_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		inAppNotification.ID,
		inAppNotification.Title,
		inAppNotification.Description,
		inAppNotification.IsRead,
		inAppNotification.ReadAt,
		inAppNotification.CreatedAt,
		inAppNotification.UpdatedAt,
	)

	if err != nil {
		return InAppNotification{}, err
	}
	
	return inAppNotification, nil
}

func (r *PostgresInAppNotificationRepository) UpdateReadAt(id string) error {
	now := time.Now()
	
	_, err := r.DB.Exec(`
		UPDATE in_app_notifications 
		SET is_read = true, read_at = $1, updated_at = $2
		WHERE id = $3
	`,
		now,
		now,
		id,
	)
	log.Println("InAppNotification READ AT updated successfully")
	return err
}

func (r *PostgresInAppNotificationRepository) Update(inAppNotification InAppNotification) error {
	inAppNotification.UpdatedAt = time.Now()
	
	_, err := r.DB.Exec(`
		UPDATE in_app_notifications 
		SET title = $1, description = $2, is_read = $3, read_at = $4, updated_at = $5
		WHERE id = $6
	`,
		inAppNotification.Title,
		inAppNotification.Description,
		inAppNotification.IsRead,
		inAppNotification.ReadAt,
		inAppNotification.UpdatedAt,
		inAppNotification.ID,
	)
	log.Println("InAppNotification updated successfully")
	return err
}

func (r *PostgresInAppNotificationRepository) Delete(id string) error {
	_, err := r.DB.Exec("DELETE FROM in_app_notifications WHERE id = $1", id)
	return err
} 