package commons

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

// @Description Task model
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // "pending", "in-progress", "completed"
	Priority    int       `json:"priority"`
	EmailSent   bool      `json:"email_sent"`
	InAppSent   bool      `json:"in_app_sent"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskRepositoryInterface interface {
	GetAll() ([]Task, error)
	GetByID(id string) (Task, error)
	Create(task Task) (Task, error)
	Update(task Task) error
	Delete(id string) error
}

type SQLiteTaskRepository struct {
	DB *sql.DB
}

func NewSQLiteTaskRepository(db *sql.DB) *SQLiteTaskRepository {
	return &SQLiteTaskRepository{DB: db}
}

func (r *SQLiteTaskRepository) GetAll() ([]Task, error) {
	rows, err := r.DB.Query(`
		SELECT id, title, description, status, priority, email_sent, in_app_sent, due_date, created_at, updated_at
		FROM tasks
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var dueDateStr, createdAtStr, updatedAtStr string
		
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.EmailSent,
			&task.InAppSent,
			&dueDateStr,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, err
		}

		task.DueDate, _ = time.Parse(time.RFC3339, dueDateStr)
		task.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		task.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)
		
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *SQLiteTaskRepository) GetByID(id string) (Task, error) {
	var task Task
	var dueDateStr, createdAtStr, updatedAtStr string
	
	err := r.DB.QueryRow(`
		SELECT id, title, description, status, priority, email_sent, in_app_sent, due_date, created_at, updated_at
		FROM tasks WHERE id = ?
	`, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.EmailSent,
		&task.InAppSent,
		&dueDateStr,
		&createdAtStr,
		&updatedAtStr,
	)
	
	if err != nil {
		return Task{}, err
	}

	task.DueDate, _ = time.Parse(time.RFC3339, dueDateStr)
	task.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	task.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)
	
	return task, nil
}

func (r *SQLiteTaskRepository) Create(task Task) (Task, error) {
	if task.ID == "" {
		task.ID = uuid.New().String()
	}

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now
	
	_, err := r.DB.Exec(`
		INSERT INTO tasks (id, title, description, status, priority, email_sent, in_app_sent, due_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.EmailSent,
		task.InAppSent,
		task.DueDate.Format(time.RFC3339),
		task.CreatedAt.Format(time.RFC3339),
		task.UpdatedAt.Format(time.RFC3339),
	)

	if err != nil {
		return Task{}, err
	}
	
	return task, nil
}

func (r *SQLiteTaskRepository) Update(task Task) error {
	task.UpdatedAt = time.Now()
	
	_, err := r.DB.Exec(`
		UPDATE tasks 
		SET title = ?, description = ?, status = ?, priority = ?, email_sent = ?, in_app_sent = ?, due_date = ?, updated_at = ?
		WHERE id = ?
	`,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.EmailSent,
		task.InAppSent,
		task.DueDate.Format(time.RFC3339),
		task.UpdatedAt.Format(time.RFC3339),
		task.ID,
	)
	log.Println("Task updated successfully")
	return err
}

func (r *SQLiteTaskRepository) Delete(id string) error {
	_, err := r.DB.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
} 
