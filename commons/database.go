package commons

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func GetPostgresConnectionString() string {
	host := GetEnv("DB_HOST", "postgres")
	port := GetEnv("DB_PORT", "5432")
	user := GetEnv("DB_USER", "postgres")
	password := GetEnv("DB_PASSWORD", "postgres")
	dbname := GetEnv("DB_NAME", "tasks")
	sslmode := GetEnv("DB_SSLMODE", "disable")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
}

func InitDB() (*sql.DB, error) {
	log.Println("Initializing PostgreSQL database...")

	connStr := GetPostgresConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Create users table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		handle VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		salt VARCHAR(32) NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'active',
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	)
	`)
	if err != nil {
		log.Printf("Error creating users table: %v", err)
		return nil, err
	}

	// Create tasks table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		creator_id TEXT NOT NULL,
		assignee_id TEXT,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL,
		priority INTEGER NOT NULL DEFAULT 0,
		email_sent BOOLEAN NOT NULL DEFAULT FALSE,
		in_app_sent BOOLEAN NOT NULL DEFAULT FALSE,
		due_date TIMESTAMP,
		deleted BOOLEAN NOT NULL DEFAULT FALSE,
		deleted_at TIMESTAMP,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		CONSTRAINT fk_tasks_creator FOREIGN KEY (creator_id)
			REFERENCES users(id) ON DELETE RESTRICT,
		CONSTRAINT fk_tasks_assignee FOREIGN KEY (assignee_id)
			REFERENCES users(id) ON DELETE SET NULL
	)
	`)
	if err != nil {
		log.Printf("Error creating tasks table: %v", err)
		return nil, err
	}

	// Create in_app_notifications table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS in_app_notifications (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		is_read BOOLEAN NOT NULL DEFAULT FALSE,
		read_at TIMESTAMP,
		deleted BOOLEAN NOT NULL DEFAULT FALSE,
		deleted_at TIMESTAMP,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		CONSTRAINT fk_notifications_user FOREIGN KEY (user_id)
			REFERENCES users(id) ON DELETE CASCADE
	)
	`)
	if err != nil {
		log.Printf("Error creating in_app_notifications table: %v", err)
		return nil, err
	}

	// Create task_system_events table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS task_system_events (
		id TEXT PRIMARY KEY,
		task_id TEXT NOT NULL,
		correlation_id TEXT NOT NULL,
		origin TEXT NOT NULL,
		action TEXT NOT NULL,	
		message TEXT NOT NULL,
		json_data TEXT,
		emit_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP NOT NULL,
		CONSTRAINT fk_task_system_events_task FOREIGN KEY (task_id)
			REFERENCES tasks(id) ON DELETE CASCADE
	)
	`)
	if err != nil {
		log.Printf("Error creating task_system_events table: %v", err)
		return nil, err
	}

	_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email)`)
	if err != nil {
		log.Printf("Warning: Failed to create unique index on users.email: %v", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_handle ON users(handle)`)
	if err != nil {
		log.Printf("Warning: Failed to create index on users.handle: %v", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)`)
	if err != nil {
		log.Printf("Warning: Failed to create index on users.status: %v", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status)`)
	if err != nil {
		log.Printf("Warning: Failed to create index on tasks.status: %v", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date)`)
	if err != nil {
		log.Printf("Warning: Failed to create index on tasks.due_date: %v", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_tasks_creator ON tasks(creator_id)`)
	if err != nil {
		log.Printf("Warning: Failed to create index on tasks.creator_id: %v", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_tasks_assignee ON tasks(assignee_id)`)
	if err != nil {
		log.Printf("Warning: Failed to create index on tasks.assignee_id: %v", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON in_app_notifications(is_read)`)
	if err != nil {
		log.Printf("Warning: Failed to create index on in_app_notifications.is_read: %v", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_notifications_user ON in_app_notifications(user_id)`)
	if err != nil {
		log.Printf("Warning: Failed to create index on in_app_notifications.user_id: %v", err)
	}

	log.Println("PostgreSQL database initialized successfully")

	return db, nil
}

func GetConnection() (*sql.DB, error) {
	connStr := GetPostgresConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("PostgreSQL database connection created")

	return db, nil
}
