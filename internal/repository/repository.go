package repository

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/auth"
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return Repository{db: db}
}

func (r Repository) User(username string) (auth.Credentials, error) {
	credentials := auth.Credentials{}
	err := r.db.QueryRow("SELECT id, hashed_password FROM users WHERE username=?", username).
		Scan(&credentials.ID, &credentials.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.Credentials{}, auth.ErrInvalidUsernameOrPassword
		}
		return auth.Credentials{}, err
	}
	return credentials, nil
}
