package commons

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

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
		var dbEvent DBTaskSystemEvent
		
		err := rows.Scan(
			&dbEvent.ID,
			&dbEvent.TaskId,
			&dbEvent.CorrelationId,
			&dbEvent.Origin,
			&dbEvent.Action,
			&dbEvent.Message,
			&dbEvent.JsonData,
			&dbEvent.EmitAt,
			&dbEvent.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		taskSystemEvents = append(taskSystemEvents, dbEvent.ToTaskSystemEvent())
	}

	return taskSystemEvents, nil
}

func (r *PostgresTaskSystemEventRepository) Create(taskSystemEvent TaskSystemEvent, delay int) (TaskSystemEvent, error) {
	dbEvent := &DBTaskSystemEvent{}
	dbEvent.FromTaskSystemEvent(taskSystemEvent)

	if dbEvent.ID == "" {
		dbEvent.ID = uuid.New().String()
	}

	dbEvent.EmitAt = time.Now().Add(time.Duration(delay) * time.Second)
	dbEvent.CreatedAt = time.Now()

	_, err := r.DB.Exec(`
		INSERT INTO task_system_events (id, task_id, correlation_id, origin, action, message, json_data, emit_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`,
		dbEvent.ID,
		dbEvent.TaskId,
		dbEvent.CorrelationId,
		dbEvent.Origin,
		dbEvent.Action,
		dbEvent.Message,
		dbEvent.JsonData,
		dbEvent.EmitAt,
		dbEvent.CreatedAt,
	)

	if err != nil {
		return TaskSystemEvent{}, err
	}
	
	return dbEvent.ToTaskSystemEvent(), nil
}
