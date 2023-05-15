package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrUserNotFound = errors.New("user not found")

type Credentials struct {
	ID           string
	PasswordHash string
}

type UserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) User(username string) (Credentials, error) {
	credentials := Credentials{}
	err := r.db.QueryRow("SELECT id, hashed_password FROM users WHERE username=?", username).
		Scan(&credentials.ID, &credentials.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Credentials{}, ErrUserNotFound
		}
		return Credentials{}, fmt.Errorf("could not get user credentials: %w", err)
	}
	return credentials, nil
}
