package repository

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/auth"
	"database/sql"
	"errors"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) User(username string) (auth.Credentials, error) {
	credentials := auth.Credentials{}
	err := r.db.QueryRow("SELECT id, hashed_password FROM users WHERE username=?", username).
		Scan(&credentials.ID, &credentials.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.Credentials{},
				fmt.Errorf("user not found: %w", auth.ErrInvalidUsernameOrPassword)
		}
		return auth.Credentials{}, fmt.Errorf("could not get user credentials: %w", err)
	}
	return credentials, nil
}
