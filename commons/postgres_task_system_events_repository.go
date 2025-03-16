package commons

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// @Description TaskSystemEvent model
type TaskSystemEvent struct {
	ID            string    `json:"id"`
	TaskId        string    `json:"task_id"`
	CorrelationId string    `json:"correlation_id"`
	Origin        string    `json:"origin"`
	Action        string    `json:"action"`
	Message       string    `json:"message"`
	JsonData      string    `json:"json_data"`
	EmitAt        time.Time `json:"emit_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type TaskSystemEventRepositoryInterface interface {
	GetAll() ([]TaskSystemEvent, error)
	Create(systemEvent TaskSystemEvent, delay int) (TaskSystemEvent, error)
}

type PostgresTaskSystemEventRepository struct {
	DB *sql.DB
}

func NewPostgresTaskSystemEventRepository(db *sql.DB) *PostgresTaskSystemEventRepository {
	return &PostgresTaskSystemEventRepository{DB: db}
}

func (r *PostgresTaskSystemEventRepository) GetAll() ([]TaskSystemEvent, error) {
	rows, err := r.DB.Query(`
		SELECT id, task_id, correlation_id, origin, action, message, json_data, emit_at, created_at
		FROM task_system_events
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taskSystemEvents []TaskSystemEvent
	for rows.Next() {
		var taskSystemEvent TaskSystemEvent
		
		err := rows.Scan(
			&taskSystemEvent.ID,
			&taskSystemEvent.TaskId,
			&taskSystemEvent.CorrelationId,
			&taskSystemEvent.Origin,
			&taskSystemEvent.Action,
			&taskSystemEvent.Message,
			&taskSystemEvent.JsonData,
			&taskSystemEvent.EmitAt,
			&taskSystemEvent.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		taskSystemEvents = append(taskSystemEvents, taskSystemEvent)
	}

	return taskSystemEvents, nil
}

func (r *PostgresTaskSystemEventRepository) Create(taskSystemEvent TaskSystemEvent, delay int) (TaskSystemEvent, error) {
	if taskSystemEvent.ID == "" {
		taskSystemEvent.ID = uuid.New().String()
	}

	taskSystemEvent.EmitAt = time.Now().Add(time.Duration(delay) * time.Second)
	taskSystemEvent.CreatedAt = time.Now()

	_, err := r.DB.Exec(`
		INSERT INTO task_system_events (id, task_id, correlation_id, origin, action, message, json_data, emit_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`,
		taskSystemEvent.ID,
		taskSystemEvent.TaskId,
		taskSystemEvent.CorrelationId,
		taskSystemEvent.Origin,
		taskSystemEvent.Action,
		taskSystemEvent.Message,
		taskSystemEvent.JsonData,
		taskSystemEvent.EmitAt,
		taskSystemEvent.CreatedAt,
	)

	if err != nil {
		return TaskSystemEvent{}, err
	}
	
	return taskSystemEvent, nil
}
