package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrUserNotFound = errors.New("user not found")

type Credentials struct {
	ID           int
	PasswordHash string
}

type Repository interface {
	GetUser(username string) (Credentials, error)
	CreateUser(username string, passwordHash string, role string) (userId int, err error)
}

type UserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) GetUser(username string) (Credentials, error) {
	credentials := Credentials{}
	err := r.db.QueryRow("SELECT user_id, hashed_password FROM users WHERE username=$1", username).
		Scan(&credentials.ID, &credentials.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Credentials{}, fmt.Errorf("%w: %v", ErrUserNotFound, err)
		}
		return Credentials{}, fmt.Errorf("could not get user credentials: %w", err)
	}
	return credentials, nil
}

func (r UserRepository) CreateUser(username string, passwordHash string, role string) (userId int, err error) {
	err = r.db.QueryRow("SELECT user_id FROM users WHERE username = $1", username).
		Scan(&userId)
	if err == nil {
		return 0, fmt.Errorf("user with username %q already exists", username)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("failed to check if user with username %q exists: %w",
			username, err)
	}

	var roleID string
	err = r.db.QueryRow("SELECT role_id FROM roles WHERE role_name = $1", role).
		Scan(&roleID)
	if err != nil {
		return 0, fmt.Errorf("could not get role id: %w", err)
	}

	err = r.db.QueryRow("INSERT INTO users (username, hashed_password, role_id) "+
		"VALUES ($1, $2, $3) RETURNING user_id", username, passwordHash, roleID).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("could not create user: %w", err)
	}

	return userId, nil
}
