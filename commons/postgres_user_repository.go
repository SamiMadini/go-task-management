package commons

import (
	"database/sql"
	"log"
	"time"
)

// @Description User model
type User struct {
	ID           string     `json:"id"`
	Handle       string     `json:"handle"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Salt         string     `json:"-"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type UserRepositoryInterface interface {
	Create(user User) (User, error)
	GetByID(id string) (User, error)
	GetByEmail(email string) (User, error)
	GetByHandle(handle string) (User, error)
}

type PostgresUserRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

func (r *PostgresUserRepository) Create(user User) (User, error) {
	log.Printf("Creating user with handle: %s", user.Handle)

	_, err := r.DB.Exec(`
		INSERT INTO users (id, handle, email, password_hash, salt, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`,
		user.ID,
		user.Handle,
		user.Email,
		user.PasswordHash,
		user.Salt,
		user.Status,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return User{}, err
	}

	log.Printf("User created successfully with ID: %s", user.ID)
	return user, nil
}

func (r *PostgresUserRepository) GetByID(id string) (User, error) {
	log.Printf("Getting user by ID: %s", id)

	var user User
	err := r.DB.QueryRow(`
		SELECT id, handle, email, password_hash, salt, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.Handle,
		&user.Email,
		&user.PasswordHash,
		&user.Salt,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		log.Printf("User not found with ID: %s", id)
		return User{}, nil
	}
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return User{}, err
	}

	log.Printf("User retrieved successfully with ID: %s", id)
	return user, nil
}

func (r *PostgresUserRepository) GetByEmail(email string) (User, error) {
	log.Printf("Getting user by email: %s", email)

	var user User
	err := r.DB.QueryRow(`
		SELECT id, handle, email, password_hash, salt, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Handle,
		&user.Email,
		&user.PasswordHash,
		&user.Salt,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		log.Printf("User not found with email: %s", email)
		return User{}, nil
	}
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return User{}, err
	}

	log.Printf("User retrieved successfully with email: %s", email)
	return user, nil
}

func (r *PostgresUserRepository) GetByHandle(handle string) (User, error) {
	log.Printf("Getting user by handle: %s", handle)

	var user User
	err := r.DB.QueryRow(`
		SELECT id, handle, email, password_hash, salt, status, created_at, updated_at
		FROM users
		WHERE handle = $1
	`, handle).Scan(
		&user.ID,
		&user.Handle,
		&user.Email,
		&user.PasswordHash,
		&user.Salt,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		log.Printf("User not found with handle: %s", handle)
		return User{}, nil
	}
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return User{}, err
	}

	log.Printf("User retrieved successfully with handle: %s", handle)
	return user, nil
}
