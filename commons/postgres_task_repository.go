package commons

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

// @Description Task model
type Task struct {
	ID          string            `json:"id"`
	CreatorID   string            `json:"creator_id"`
	AssigneeID  *string           `json:"assignee_id,omitempty"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      string            `json:"status"` // "todo", "in-progress", "done"
	Priority    int               `json:"priority"`
	EmailSent   bool              `json:"email_sent"`
	InAppSent   bool              `json:"in_app_sent"`
	DueDate     time.Time         `json:"due_date"`
	Deleted     bool              `json:"deleted"`
	DeletedAt   *time.Time        `json:"deleted_at,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Events      []TaskSystemEvent `json:"events"`
}

type TaskRepositoryInterface interface {
	GetAll() ([]Task, error)
	GetByID(id string) (Task, error)
	GetByUserID(userID string) ([]Task, error)
	Create(task Task) (Task, error)
	Update(task Task) error
	Delete(id string) error
	HardDelete(id string) error
}

type PostgresTaskRepository struct {
	DB *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{DB: db}
}

func (r *PostgresTaskRepository) GetAll() ([]Task, error) {
	log.Println("Getting all tasks with related events")

	// Single query with LEFT JOIN to get tasks with their events
	rows, err := r.DB.Query(`
		SELECT 
			t.id, t.creator_id, t.assignee_id, t.title, t.description, t.status, t.priority,
			t.email_sent, t.in_app_sent, t.due_date, t.created_at, t.updated_at, t.deleted, t.deleted_at,
			e.id, e.task_id, e.correlation_id, e.origin, e.action,
			e.message, e.json_data, e.emit_at, e.created_at
		FROM tasks t
		WHERE t.deleted = false
		LEFT JOIN task_system_events e ON t.id = e.task_id
		ORDER BY t.created_at DESC, e.created_at DESC
	`)
	if err != nil {
		log.Printf("Failed to get tasks with events: %v", err)
		return nil, err
	}
	defer rows.Close()

	tasksMap := make(map[string]*Task)
	var tasks []Task

	for rows.Next() {
		var task Task
		var dueDate sql.NullTime
		var assigneeID sql.NullString

		// For event fields - all are nullable because of LEFT JOIN
		var eventID, eventTaskId, eventCorrelationId, eventOrigin, eventAction, eventMessage, eventJsonData sql.NullString
		var eventEmitAt, eventCreatedAt sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.CreatorID,
			&assigneeID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.EmailSent,
			&task.InAppSent,
			&dueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.Deleted,
			&task.DeletedAt,
			&eventID,
			&eventTaskId,
			&eventCorrelationId,
			&eventOrigin,
			&eventAction,
			&eventMessage,
			&eventJsonData,
			&eventEmitAt,
			&eventCreatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan task and event: %v", err)
			return nil, err
		}

		if dueDate.Valid {
			task.DueDate = dueDate.Time
		}

		if assigneeID.Valid {
			task.AssigneeID = &assigneeID.String
		}

		existingTask, exists := tasksMap[task.ID]
		if !exists {
			task.Events = []TaskSystemEvent{}
			tasks = append(tasks, task)
			tasksMap[task.ID] = &tasks[len(tasks)-1]
			existingTask = &tasks[len(tasks)-1]
		}

		if eventID.Valid {
			event := TaskSystemEvent{
				ID:            eventID.String,
				TaskId:        eventTaskId.String,
				CorrelationId: eventCorrelationId.String,
				Origin:        eventOrigin.String,
				Action:        eventAction.String,
				Message:       eventMessage.String,
				JsonData:      eventJsonData.String,
				EmitAt:        eventEmitAt.Time,
				CreatedAt:     eventCreatedAt.Time,
			}
			existingTask.Events = append(existingTask.Events, event)
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	log.Println("Tasks with events retrieved successfully")
	return tasks, nil
}

func (r *PostgresTaskRepository) GetByID(id string) (Task, error) {
	log.Println("Getting task with related events by ID:", id)

	// Single query with LEFT JOIN to get task with events
	rows, err := r.DB.Query(`
		SELECT 
			t.id, t.creator_id, t.assignee_id, t.title, t.description, t.status, t.priority,
			t.email_sent, t.in_app_sent, t.due_date, t.created_at, t.updated_at, t.deleted, t.deleted_at,
			e.id, e.task_id, e.correlation_id, e.origin, e.action, e.message, e.json_data, e.emit_at, e.created_at
		FROM tasks t
		LEFT JOIN task_system_events e ON t.id = e.task_id
		WHERE t.id = $1
		ORDER BY e.created_at DESC
	`, id)
	
	if err != nil {
		log.Printf("Failed to get task with events: %v", err)
		return Task{}, err
	}
	defer rows.Close()

	var task Task
	found := false

	for rows.Next() {
		var dueDate sql.NullTime
		var assigneeID sql.NullString

		var eventID, eventTaskID, eventCorrelationID, eventOrigin, eventAction, eventMessage, eventJsonData sql.NullString
		var eventEmitAt, eventCreatedAt sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.CreatorID,
			&assigneeID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.EmailSent,
			&task.InAppSent,
			&dueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.Deleted,
			&task.DeletedAt,
			&eventID,
			&eventTaskID,
			&eventCorrelationID,
			&eventOrigin,
			&eventAction,
			&eventMessage,
			&eventJsonData,
			&eventEmitAt,
			&eventCreatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan task and event: %v", err)
			return Task{}, err
		}

		if dueDate.Valid {
			task.DueDate = dueDate.Time
		}

		if assigneeID.Valid {
			task.AssigneeID = &assigneeID.String
		}

		if !found {
			task.Events = []TaskSystemEvent{}
			found = true
		}

		if eventID.Valid && eventTaskID.Valid {
			event := TaskSystemEvent{
				ID:            eventID.String,
				TaskId:        eventTaskID.String,
				CorrelationId: eventCorrelationID.String,
				Origin:        eventOrigin.String,
				Action:        eventAction.String,
				Message:       eventMessage.String,
				JsonData:      eventJsonData.String,
			}

			if eventEmitAt.Valid {
				event.EmitAt = eventEmitAt.Time
			}

			if eventCreatedAt.Valid {
				event.CreatedAt = eventCreatedAt.Time
			}

			task.Events = append(task.Events, event)
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return Task{}, err
	}

	if !found {
		return Task{}, sql.ErrNoRows
	}

	log.Println("Task with events retrieved successfully")
	return task, nil
}

func (r *PostgresTaskRepository) GetByUserID(userID string) ([]Task, error) {
	rows, err := r.DB.Query(`
		SELECT t.id, t.creator_id, t.assignee_id, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at,
			json_agg(json_build_object(
				'id', e.id,
				'task_id', e.task_id,
				'correlation_id', e.correlation_id,
				'origin', e.origin,
				'action', e.action,
				'message', e.message,
				'json_data', e.json_data,
				'emit_at', to_char(e.emit_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
				'created_at', to_char(e.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
			)) as events
		FROM tasks t
		LEFT JOIN task_system_events e ON t.id = e.task_id
		WHERE (t.creator_id = $1 OR t.assignee_id = $1) AND t.deleted = false
		GROUP BY t.id, t.creator_id, t.assignee_id, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at
		ORDER BY t.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var eventsJSON []byte
		var assigneeID sql.NullString
		var dueDate sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.CreatorID,
			&assigneeID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&dueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
			&eventsJSON,
		)
		if err != nil {
			return nil, err
		}

		if assigneeID.Valid {
			task.AssigneeID = &assigneeID.String
		}

		if dueDate.Valid {
			task.DueDate = dueDate.Time
		}

		if eventsJSON != nil {
			if err := json.Unmarshal(eventsJSON, &task.Events); err != nil {
				return nil, err
			}
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *PostgresTaskRepository) Create(task Task) (Task, error) {
	_, err := r.DB.Exec(`
		INSERT INTO tasks (id, creator_id, assignee_id, title, description, status, priority, email_sent, in_app_sent, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`,
		task.ID,
		task.CreatorID,
		task.AssigneeID,
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
		SET title = $1, description = $2, status = $3, priority = $4, email_sent = $5, in_app_sent = $6, due_date = $7, assignee_id = $8, updated_at = $9
		WHERE id = $10
	`,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.EmailSent,
		task.InAppSent,
		task.DueDate,
		task.AssigneeID,
		task.UpdatedAt,
		task.ID,
	)
	log.Println("Task updated successfully")
	return err
}

func (r *PostgresTaskRepository) Delete(id string) error {
	_, err := r.DB.Exec(`
		UPDATE tasks 
		SET deleted = $1, deleted_at = $2, updated_at = $3
		WHERE id = $4
	`,
		true,
		time.Now(),
		time.Now(),
		id,
	)
	log.Println("Task soft deleted successfully")
	return err
}

func (r *PostgresTaskRepository) HardDelete(id string) error {
	_, err := r.DB.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		log.Printf("Error hard deleting task with ID %s: %v", id, err)
		var exists bool
		err = r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1)", id).Scan(&exists)
		if err != nil {
			log.Printf("Error checking if task exists: %v", err)
		} else if !exists {
			log.Printf("Task with ID %s does not exist", id)
		}
	} else {
		log.Printf("Task with ID %s hard deleted successfully", id)
	}
	return err
}
