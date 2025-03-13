package commons

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	log.Println("Initializing database...")

	var dbPath = GetEnv("DB_PATH", "../tasks.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL,
		priority INTEGER NOT NULL DEFAULT 0,
		email_sent BOOLEAN NOT NULL DEFAULT FALSE,
		in_app_sent BOOLEAN NOT NULL DEFAULT FALSE,
		due_date TEXT,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	)
	`)
	if err != nil {
		log.Printf("Error creating tasks table: %v", err)
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS in_app_notifications (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		is_read BOOLEAN NOT NULL DEFAULT FALSE,
		read_at TEXT DEFAULT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	)
	`)
	if err != nil {
		log.Printf("Error creating in_app_notifications table: %v", err)
		return nil, err
	}

	log.Println("Database initialized successfully")

	return db, nil
}

func GetConnection() (*sql.DB, error) {
	var dbPath = GetEnv("DB_PATH", "../tasks.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Database connection created")

	return db, nil
}