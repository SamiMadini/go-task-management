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

type SQLiteInAppNotificationRepository struct {
	DB *sql.DB
}

func NewSQLiteInAppNotificationRepository(db *sql.DB) *SQLiteInAppNotificationRepository {
	return &SQLiteInAppNotificationRepository{DB: db}
}

func (r *SQLiteInAppNotificationRepository) GetAll() ([]InAppNotification, error) {
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
		var createdAtStr, updatedAtStr string
		
		err := rows.Scan(
			&inAppNotification.ID,
			&inAppNotification.Title,
			&inAppNotification.Description,
			&inAppNotification.IsRead,
			&inAppNotification.ReadAt,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, err
		}

		inAppNotification.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		inAppNotification.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)
		
		inAppNotifications = append(inAppNotifications, inAppNotification)
	}

	return inAppNotifications, nil
}

func (r *SQLiteInAppNotificationRepository) GetByID(id string) (InAppNotification, error) {
	var inAppNotification InAppNotification
	var createdAtStr, updatedAtStr string
	
	err := r.DB.QueryRow(`
		SELECT id, title, description, is_read, read_at, created_at, updated_at
		FROM in_app_notifications WHERE id = ?
	`, id).Scan(
		&inAppNotification.ID,
		&inAppNotification.Title,
		&inAppNotification.Description,
		&inAppNotification.IsRead,
		&inAppNotification.ReadAt,
		&createdAtStr,
		&updatedAtStr,
	)
	
	if err != nil {
		return InAppNotification{}, err
	}

	inAppNotification.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	inAppNotification.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)
	
	return inAppNotification, nil
}

func (r *SQLiteInAppNotificationRepository) Create(inAppNotification InAppNotification) (InAppNotification, error) {
	if inAppNotification.ID == "" {
		inAppNotification.ID = uuid.New().String()
	}

	now := time.Now()
	inAppNotification.CreatedAt = now
	inAppNotification.UpdatedAt = now
	
	_, err := r.DB.Exec(`
		INSERT INTO in_app_notifications (id, title, description, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?)
	`,
		inAppNotification.ID,
		inAppNotification.Title,
		inAppNotification.Description,
		inAppNotification.CreatedAt.Format(time.RFC3339),
		inAppNotification.UpdatedAt.Format(time.RFC3339),
	)

	if err != nil {
		return InAppNotification{}, err
	}
	
	return inAppNotification, nil
}

func (r *SQLiteInAppNotificationRepository) UpdateReadAt(id string) error {
	_, err := r.DB.Exec(`
		UPDATE in_app_notifications 
		SET read_at = ?, updated_at = ?
		WHERE id = ?
	`,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339),
		id,
	)
	log.Println("InAppNotification READ AT updated successfully")
	return err
}

func (r *SQLiteInAppNotificationRepository) Update(inAppNotification InAppNotification) error {
	inAppNotification.UpdatedAt = time.Now()
	
	_, err := r.DB.Exec(`
		UPDATE in_app_notifications 
		SET title = ?, description = ?, is_read = ?, read_at = ?, updated_at = ?
		WHERE id = ?
	`,
		inAppNotification.Title,
		inAppNotification.Description,
		inAppNotification.IsRead,
		inAppNotification.ReadAt,
		inAppNotification.UpdatedAt.Format(time.RFC3339),
		inAppNotification.ID,
	)
	log.Println("InAppNotification updated successfully")
	return err
}

func (r *SQLiteInAppNotificationRepository) Delete(id string) error {
	_, err := r.DB.Exec("DELETE FROM in_app_notifications WHERE id = ?", id)
	return err
}

