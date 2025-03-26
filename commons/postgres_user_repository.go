package commons

import (
	"database/sql"
	"log"
	"time"
)

type UserRepositoryInterface interface {
	Create(user User) (User, error)
	GetByID(id string) (User, error)
	GetByEmail(email string) (User, error)
	GetByHandle(handle string) (User, error)
	UpdatePassword(id string, hashedPassword string, salt string) (User, error)
}

type PostgresUserRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

func (r *PostgresUserRepository) Create(user User) (User, error) {
	log.Printf("Creating user with handle: %s", user.Handle)

	dbUser := &DBUser{}
	dbUser.FromUser(user)
	dbUser.CreatedAt = time.Now()
	dbUser.UpdatedAt = time.Now()

	_, err := r.DB.Exec(`
		INSERT INTO users (id, handle, email, password_hash, salt, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`,
		dbUser.ID,
		dbUser.Handle,
		dbUser.Email,
		dbUser.HashedPassword,
		dbUser.Salt,
		dbUser.Status,
		dbUser.CreatedAt,
		dbUser.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return User{}, err
	}

	log.Printf("User created successfully with ID: %s", dbUser.ID)
	return dbUser.ToUser(), nil
}

func (r *PostgresUserRepository) GetByID(id string) (User, error) {
	log.Printf("Getting user by ID: %s", id)

	var dbUser DBUser
	err := r.DB.QueryRow(`
		SELECT id, handle, email, password_hash, salt, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&dbUser.ID,
		&dbUser.Handle,
		&dbUser.Email,
		&dbUser.HashedPassword,
		&dbUser.Salt,
		&dbUser.Status,
		&dbUser.CreatedAt,
		&dbUser.UpdatedAt,
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
	return dbUser.ToUser(), nil
}

func (r *PostgresUserRepository) GetByEmail(email string) (User, error) {
	log.Printf("Getting user by email: %s", email)

	var dbUser DBUser
	err := r.DB.QueryRow(`
		SELECT id, handle, email, password_hash, salt, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&dbUser.ID,
		&dbUser.Handle,
		&dbUser.Email,
		&dbUser.HashedPassword,
		&dbUser.Salt,
		&dbUser.Status,
		&dbUser.CreatedAt,
		&dbUser.UpdatedAt,
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
	return dbUser.ToUser(), nil
}

func (r *PostgresUserRepository) GetByHandle(handle string) (User, error) {
	log.Printf("Getting user by handle: %s", handle)

	var dbUser DBUser
	err := r.DB.QueryRow(`
		SELECT id, handle, email, password_hash, salt, status, created_at, updated_at
		FROM users
		WHERE handle = $1
	`, handle).Scan(
		&dbUser.ID,
		&dbUser.Handle,
		&dbUser.Email,
		&dbUser.HashedPassword,
		&dbUser.Salt,
		&dbUser.Status,
		&dbUser.CreatedAt,
		&dbUser.UpdatedAt,
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
	return dbUser.ToUser(), nil
}

func (r *PostgresUserRepository) UpdatePassword(id string, hashedPassword string, salt string) (User, error) {
	var dbUser DBUser
	query := `
		UPDATE users 
		SET password_hash = $1, salt = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, handle, email, password_hash, salt, status, created_at, updated_at
	`
	err := r.DB.QueryRow(query, hashedPassword, salt, id).Scan(
		&dbUser.ID,
		&dbUser.Handle,
		&dbUser.Email,
		&dbUser.HashedPassword,
		&dbUser.Salt,
		&dbUser.Status,
		&dbUser.CreatedAt,
		&dbUser.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error updating user password: %v", err)
		return User{}, err
	}
	return dbUser.ToUser(), nil
}
