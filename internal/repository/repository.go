package repository

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/auth"
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return Repository{db: db}
}

func (r Repository) User(username string) (auth.Credentials, error) {
	credentials := auth.Credentials{}
	err := r.db.QueryRow("SELECT user_id, hashed_password FROM users WHERE username=$1", username).
		Scan(&credentials.ID, &credentials.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.Credentials{},
				fmt.Errorf("failed to find user in database: %w", auth.ErrInvalidUsernameOrPassword)
		}
		return auth.Credentials{}, fmt.Errorf("failed to get user credentials from database: %w", err)
	}

	return credentials, nil
}

func (r Repository) CreateUser(username string, passwordHash string, role string) (string, error) {
	var id string
	err := r.db.QueryRow("SELECT user_id FROM users WHERE username = $1", username).Scan(&id)
	if err == nil {
		return "", fmt.Errorf("user with username %q already exists", username)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("failed to check if user with username %q exists: %w", username, err)
	}

	var roleID string
	err = r.db.QueryRow("SELECT role_id FROM roles WHERE role_name = $1", role).Scan(&roleID)
	if err != nil {
		return "", fmt.Errorf("failed to get role id from database: %w", err)
	}

	err = r.db.QueryRow("INSERT INTO users (username, hashed_password, role_id) VALUES ($1, $2, $3) "+
		"RETURNING user_id", username, passwordHash, roleID).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to insert new user into database: %w", err)
	}

	return id, nil
}
