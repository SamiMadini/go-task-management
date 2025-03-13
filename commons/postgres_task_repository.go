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

type PostgresTaskRepository struct {
	DB *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{DB: db}
}

func (r *PostgresTaskRepository) GetAll() ([]Task, error) {
	log.Println("Getting all tasks")

	rows, err := r.DB.Query(`
		SELECT id, title, description, status, priority, email_sent, in_app_sent, due_date, created_at, updated_at
		FROM tasks
	`)
	if err != nil {
		log.Printf("Failed to get all tasks: %v", err)
		return nil, err
	}
	defer rows.Close()

	log.Println("Go returns the task")

	var tasks []Task
	for rows.Next() {
		var task Task
		var dueDate sql.NullTime
		
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.EmailSent,
			&task.InAppSent,
			&dueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to SCAN task: %v", err)
			return nil, err
		}

		if dueDate.Valid {
			task.DueDate = dueDate.Time
		}
		
		tasks = append(tasks, task)
	}

	log.Println("Tasks: ", tasks)
	return tasks, nil
}

func (r *PostgresTaskRepository) GetByID(id string) (Task, error) {
	var task Task
	var dueDate sql.NullTime
	
	err := r.DB.QueryRow(`
		SELECT id, title, description, status, priority, email_sent, in_app_sent, due_date, created_at, updated_at
		FROM tasks WHERE id = $1
	`, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.EmailSent,
		&task.InAppSent,
		&dueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	
	if err != nil {
		return Task{}, err
	}

	if dueDate.Valid {
		task.DueDate = dueDate.Time
	}
	
	return task, nil
}

func (r *PostgresTaskRepository) Create(task Task) (Task, error) {
	if task.ID == "" {
		task.ID = uuid.New().String()
	}

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now
	
	_, err := r.DB.Exec(`
		INSERT INTO tasks (id, title, description, status, priority, email_sent, in_app_sent, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.EmailSent,
		task.InAppSent,
		task.DueDate,
		task.CreatedAt,
		task.UpdatedAt,
	)

	if err != nil {
		return Task{}, err
	}
	
	return task, nil
}

func (r *PostgresTaskRepository) Update(task Task) error {
	task.UpdatedAt = time.Now()
	
	_, err := r.DB.Exec(`
		UPDATE tasks 
		SET title = $1, description = $2, status = $3, priority = $4, email_sent = $5, in_app_sent = $6, due_date = $7, updated_at = $8
		WHERE id = $9
	`,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.EmailSent,
		task.InAppSent,
		task.DueDate,
		task.UpdatedAt,
		task.ID,
	)
	log.Println("Task updated successfully")
	return err
}

func (r *PostgresTaskRepository) Delete(id string) error {
	_, err := r.DB.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting task with ID %s: %v", id, err)
		// Check if the task exists before trying to delete
		var exists bool
		err = r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1)", id).Scan(&exists)
		if err != nil {
			log.Printf("Error checking if task exists: %v", err)
		} else if !exists {
			log.Printf("Task with ID %s does not exist", id)
		}
	} else {
		log.Printf("Task with ID %s deleted successfully", id)
	}
	return err
} 