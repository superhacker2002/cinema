package repository

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/auth"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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
			return auth.Credentials{}, auth.ErrInvalidUsernameOrPassword
		}
		return auth.Credentials{}, err
	}
	return credentials, nil
}

func (r Repository) CreateUser(username string, password string) (string, error) {
	var id string
	err := r.db.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&id)
	if err == nil {
		return "", fmt.Errorf("user with username %q already exists", username)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("failed to check if user with username %q exists: %w", username, err)
	}

	hash := sha256.Sum256([]byte(password))
	hashHex := hex.EncodeToString(hash[:])

	err = r.db.QueryRow("INSERT INTO users (username, hashed_password) VALUES ($1, $2) RETURNING id",
		username, hashHex).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to insert new user into database: %w", err)
	}

	return id, nil
}
