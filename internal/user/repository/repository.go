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
}

type UserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) GetUser(username string) (Credentials, error) {
	credentials := Credentials{}
	err := r.db.QueryRow("SELECT id, hashed_password FROM users WHERE username=?", username).
		Scan(&credentials.ID, &credentials.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Credentials{}, fmt.Errorf("%w: %v", ErrUserNotFound, err)
		}
		return Credentials{}, fmt.Errorf("could not get user credentials: %w", err)
	}
	return credentials, nil
}
