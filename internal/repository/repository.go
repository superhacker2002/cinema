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

func (r repository) getUserInfo(username string) (auth.Credentials, error) {
	credentials := auth.Credentials{}
	err := r.db.QueryRow("SELECT id, password_hash FROM users WHERE username=?", username).
		Scan(&credentials.ID, &credentials.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.Credentials{}, errors.New("user not found")
		}
		return auth.Credentials{}, err
	}
	return credentials, nil
}
