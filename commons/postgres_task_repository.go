package commons

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

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

	rows, err := r.DB.Query(`
		SELECT 
			t.id, t.creator_id, t.assignee_id, t.title, t.description, t.status, t.priority,
			t.email_sent, t.in_app_sent, t.due_date, t.created_at, t.updated_at, t.deleted, t.deleted_at,
			e.id, e.task_id, e.correlation_id, e.origin, e.action,
			e.message, e.json_data, e.emit_at, e.created_at
		FROM tasks t
		LEFT JOIN task_system_events e ON t.id = e.task_id
		WHERE t.deleted = false
		ORDER BY t.created_at DESC, e.created_at DESC
	`)
	if err != nil {
		log.Printf("Failed to get tasks with events: %v", err)
		return nil, err
	}
	defer rows.Close()

	tasksMap := make(map[string]*DBTask)
	var tasks []Task

	for rows.Next() {
		var dbTask DBTask
		var dueDate sql.NullTime
		var assigneeID sql.NullString

		var eventID, eventTaskId, eventCorrelationId, eventOrigin, eventAction, eventMessage, eventJsonData sql.NullString
		var eventEmitAt, eventCreatedAt sql.NullTime

		err := rows.Scan(
			&dbTask.ID,
			&dbTask.CreatorID,
			&assigneeID,
			&dbTask.Title,
			&dbTask.Description,
			&dbTask.Status,
			&dbTask.Priority,
			&dbTask.EmailSent,
			&dbTask.InAppSent,
			&dueDate,
			&dbTask.CreatedAt,
			&dbTask.UpdatedAt,
			&dbTask.Deleted,
			&dbTask.DeletedAt,
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
			dbTask.DueDate = dueDate.Time
		}

		if assigneeID.Valid {
			dbTask.AssigneeID = &assigneeID.String
		}

		existingTask, exists := tasksMap[dbTask.ID]
		if !exists {
			dbTask.Events = []TaskSystemEvent{}
			tasksMap[dbTask.ID] = &dbTask
			existingTask = &dbTask
		}

		if eventID.Valid {
			dbEvent := DBTaskSystemEvent{
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
			existingTask.Events = append(existingTask.Events, dbEvent.ToTaskSystemEvent())
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	// Convert DB tasks to domain tasks
	for _, dbTask := range tasksMap {
		tasks = append(tasks, dbTask.ToTask())
	}

	log.Println("Tasks with events retrieved successfully")
	return tasks, nil
}

func (r *PostgresTaskRepository) GetByID(id string) (Task, error) {
	log.Println("Getting task with related events by ID:", id)

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

	var dbTask DBTask
	found := false

	for rows.Next() {
		var dueDate sql.NullTime
		var assigneeID sql.NullString

		var eventID, eventTaskID, eventCorrelationID, eventOrigin, eventAction, eventMessage, eventJsonData sql.NullString
		var eventEmitAt, eventCreatedAt sql.NullTime

		err := rows.Scan(
			&dbTask.ID,
			&dbTask.CreatorID,
			&assigneeID,
			&dbTask.Title,
			&dbTask.Description,
			&dbTask.Status,
			&dbTask.Priority,
			&dbTask.EmailSent,
			&dbTask.InAppSent,
			&dueDate,
			&dbTask.CreatedAt,
			&dbTask.UpdatedAt,
			&dbTask.Deleted,
			&dbTask.DeletedAt,
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
			dbTask.DueDate = dueDate.Time
		}

		if assigneeID.Valid {
			dbTask.AssigneeID = &assigneeID.String
		}

		if !found {
			dbTask.Events = []TaskSystemEvent{}
			found = true
		}

		if eventID.Valid && eventTaskID.Valid {
			dbEvent := DBTaskSystemEvent{
				ID:            eventID.String,
				TaskId:        eventTaskID.String,
				CorrelationId: eventCorrelationID.String,
				Origin:        eventOrigin.String,
				Action:        eventAction.String,
				Message:       eventMessage.String,
				JsonData:      eventJsonData.String,
			}

			if eventEmitAt.Valid {
				dbEvent.EmitAt = eventEmitAt.Time
			}

			if eventCreatedAt.Valid {
				dbEvent.CreatedAt = eventCreatedAt.Time
			}

			dbTask.Events = append(dbTask.Events, dbEvent.ToTaskSystemEvent())
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
	return dbTask.ToTask(), nil
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
		var dbTask DBTask
		var eventsJSON []byte
		var assigneeID sql.NullString
		var dueDate sql.NullTime

		err := rows.Scan(
			&dbTask.ID,
			&dbTask.CreatorID,
			&assigneeID,
			&dbTask.Title,
			&dbTask.Description,
			&dbTask.Status,
			&dbTask.Priority,
			&dueDate,
			&dbTask.CreatedAt,
			&dbTask.UpdatedAt,
			&eventsJSON,
		)
		if err != nil {
			return nil, err
		}

		if assigneeID.Valid {
			dbTask.AssigneeID = &assigneeID.String
		}

		if dueDate.Valid {
			dbTask.DueDate = dueDate.Time
		}

		if eventsJSON != nil {
			var dbEvents []DBTaskSystemEvent
			if err := json.Unmarshal(eventsJSON, &dbEvents); err != nil {
				return nil, err
			}
			for _, dbEvent := range dbEvents {
				dbTask.Events = append(dbTask.Events, dbEvent.ToTaskSystemEvent())
			}
		}

		tasks = append(tasks, dbTask.ToTask())
	}

	return tasks, nil
}

func (r *PostgresTaskRepository) Create(task Task) (Task, error) {
	dbTask := &DBTask{}
	dbTask.FromTask(task)
	dbTask.CreatedAt = time.Now()
	dbTask.UpdatedAt = time.Now()

	_, err := r.DB.Exec(`
		INSERT INTO tasks (id, creator_id, assignee_id, title, description, status, priority, email_sent, in_app_sent, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`,
		dbTask.ID,
		dbTask.CreatorID,
		dbTask.AssigneeID,
		dbTask.Title,
		dbTask.Description,
		dbTask.Status,
		dbTask.Priority,
		dbTask.EmailSent,
		dbTask.InAppSent,
		dbTask.DueDate,
		dbTask.CreatedAt,
		dbTask.UpdatedAt,
	)

	if err != nil {
		return Task{}, err
	}

	return dbTask.ToTask(), nil
}

func (r *PostgresTaskRepository) Update(task Task) error {
	dbTask := &DBTask{}
	dbTask.FromTask(task)
	dbTask.UpdatedAt = time.Now()
	
	_, err := r.DB.Exec(`
		UPDATE tasks 
		SET title = $1, description = $2, status = $3, priority = $4, email_sent = $5, in_app_sent = $6, due_date = $7, assignee_id = $8, updated_at = $9
		WHERE id = $10
	`,
		dbTask.Title,
		dbTask.Description,
		dbTask.Status,
		dbTask.Priority,
		dbTask.EmailSent,
		dbTask.InAppSent,
		dbTask.DueDate,
		dbTask.AssigneeID,
		dbTask.UpdatedAt,
		dbTask.ID,
	)
	log.Println("Task updated successfully")
	return err
}

func (r *PostgresTaskRepository) Delete(id string) error {
	now := time.Now()
	_, err := r.DB.Exec(`
		UPDATE tasks 
		SET deleted = $1, deleted_at = $2, updated_at = $3
		WHERE id = $4
	`,
		true,
		now,
		now,
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
