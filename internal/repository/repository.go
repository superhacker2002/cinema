package repository

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/auth"
	"database/sql"
	"errors"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) repository {
	return repository{db: db}
}

func (r repository) GetUserInfo(username string) (auth.Credentials, error) {
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
