package commons

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type InAppNotificationRepositoryInterface interface {
	GetAll() ([]InAppNotification, error)
	GetByID(id string) (InAppNotification, error)
	GetByUserID(userID string) ([]InAppNotification, error)
	Create(inAppNotification InAppNotification) (InAppNotification, error)
	UpdateOnRead(id string, isRead bool) error
	Update(inAppNotification InAppNotification) error
	Delete(id string) error
	HardDelete(id string) error
}

type PostgresInAppNotificationRepository struct {
	DB *sql.DB
}

func NewPostgresInAppNotificationRepository(db *sql.DB) *PostgresInAppNotificationRepository {
	return &PostgresInAppNotificationRepository{DB: db}
}

func (r *PostgresInAppNotificationRepository) GetAll() ([]InAppNotification, error) {
	rows, err := r.DB.Query(`
		SELECT id, user_id, title, description, is_read, read_at, created_at, updated_at, deleted, deleted_at
		FROM in_app_notifications
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inAppNotifications []InAppNotification
	for rows.Next() {
		var dbNotification DBInAppNotification
		var readAt sql.NullTime

		err := rows.Scan(
			&dbNotification.ID,
			&dbNotification.UserID,
			&dbNotification.Title,
			&dbNotification.Description,
			&dbNotification.IsRead,
			&readAt,
			&dbNotification.CreatedAt,
			&dbNotification.UpdatedAt,
			&dbNotification.Deleted,
			&dbNotification.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		if readAt.Valid {
			dbNotification.ReadAt = &readAt.Time
		}

		inAppNotifications = append(inAppNotifications, dbNotification.ToInAppNotification())
	}

	return inAppNotifications, nil
}

func (r *PostgresInAppNotificationRepository) GetByID(id string) (InAppNotification, error) {
	var dbNotification DBInAppNotification
	var readAt sql.NullTime

	err := r.DB.QueryRow(`
		SELECT id, user_id, title, description, is_read, read_at, created_at, updated_at, deleted, deleted_at
		FROM in_app_notifications WHERE id = $1
	`, id).Scan(
		&dbNotification.ID,
		&dbNotification.UserID,
		&dbNotification.Title,
		&dbNotification.Description,
		&dbNotification.IsRead,
		&readAt,
		&dbNotification.CreatedAt,
		&dbNotification.UpdatedAt,
		&dbNotification.Deleted,
		&dbNotification.DeletedAt,
	)
	
	if err != nil {
		return InAppNotification{}, err
	}

	if readAt.Valid {
		dbNotification.ReadAt = &readAt.Time
	}

	return dbNotification.ToInAppNotification(), nil
}

func (r *PostgresInAppNotificationRepository) GetByUserID(userID string) ([]InAppNotification, error) {
	rows, err := r.DB.Query(`
		SELECT id, user_id, title, description, is_read, read_at, created_at, updated_at, deleted, deleted_at
		FROM in_app_notifications
		WHERE user_id = $1 AND deleted = false
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inAppNotifications []InAppNotification
	for rows.Next() {
		var dbNotification DBInAppNotification
		var readAt sql.NullTime

		err := rows.Scan(
			&dbNotification.ID,
			&dbNotification.UserID,
			&dbNotification.Title,
			&dbNotification.Description,
			&dbNotification.IsRead,
			&readAt,
			&dbNotification.CreatedAt,
			&dbNotification.UpdatedAt,
			&dbNotification.Deleted,
			&dbNotification.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		if readAt.Valid {
			dbNotification.ReadAt = &readAt.Time
		}

		inAppNotifications = append(inAppNotifications, dbNotification.ToInAppNotification())
	}

	return inAppNotifications, nil
}

func (r *PostgresInAppNotificationRepository) Create(notification InAppNotification) (InAppNotification, error) {
	dbNotification := &DBInAppNotification{}
	dbNotification.FromInAppNotification(notification)

	if dbNotification.ID == "" {
		dbNotification.ID = uuid.New().String()
	}

	now := time.Now()
	dbNotification.CreatedAt = now
	dbNotification.UpdatedAt = now

	_, err := r.DB.Exec(`
		INSERT INTO in_app_notifications (id, user_id, title, description, is_read, read_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`,
		dbNotification.ID,
		dbNotification.UserID,
		dbNotification.Title,
		dbNotification.Description,
		dbNotification.IsRead,
		dbNotification.ReadAt,
		dbNotification.CreatedAt,
		dbNotification.UpdatedAt,
	)

	if err != nil {
		return InAppNotification{}, err
	}
	
	return dbNotification.ToInAppNotification(), nil
}

func (r *PostgresInAppNotificationRepository) UpdateOnRead(id string, isRead bool) error {
	now := time.Now()

	var readAt *time.Time
	if isRead {
		readAt = &now
	} else {
		readAt = nil
	}

	_, err := r.DB.Exec(`
		UPDATE in_app_notifications 
		SET is_read = $1, read_at = $2, updated_at = $3
		WHERE id = $4
	`,
		isRead,
		readAt,
		now,
		id,
	)
	log.Println("InAppNotification updated ON READ successfully")
	return err
}

func (r *PostgresInAppNotificationRepository) Update(notification InAppNotification) error {
	dbNotification := &DBInAppNotification{}
	dbNotification.FromInAppNotification(notification)
	dbNotification.UpdatedAt = time.Now()
	
	_, err := r.DB.Exec(`
		UPDATE in_app_notifications 
		SET title = $1, description = $2, is_read = $3, read_at = $4, updated_at = $5
		WHERE id = $6
	`,
		dbNotification.Title,
		dbNotification.Description,
		dbNotification.IsRead,
		dbNotification.ReadAt,
		dbNotification.UpdatedAt,
		dbNotification.ID,
	)
	log.Println("InAppNotification updated successfully")
	return err
}

func (r *PostgresInAppNotificationRepository) Delete(id string) error {
	_, err := r.DB.Exec(`
		UPDATE in_app_notifications 
		SET deleted = $1, deleted_at = $2, updated_at = $3 
		WHERE id = $4
	`,
		true,
		time.Now(),
		time.Now(),
		id,
	)
	log.Println("InAppNotification soft deleted successfully")
	return err
}

func (r *PostgresInAppNotificationRepository) HardDelete(id string) error {
	_, err := r.DB.Exec("DELETE FROM in_app_notifications WHERE id = $1", id)
	return err
}
