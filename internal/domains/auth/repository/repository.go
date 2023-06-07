package repository

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/auth/service"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type AuthRepository struct {
	db *sql.DB
}

func New(db *sql.DB) AuthRepository {
	return AuthRepository{db: db}
}

func (a AuthRepository) GetUser(username string) (service.Credentials, error) {
	credentials := service.Credentials{}
	err := a.db.QueryRow("SELECT user_id, hashed_password FROM users WHERE username=$1", username).
		Scan(&credentials.ID, &credentials.PasswordHash)

	if errors.Is(err, sql.ErrNoRows) {
		return service.Credentials{}, service.ErrUserNotFound
	}

	if err != nil {
		return service.Credentials{}, fmt.Errorf("could not get user credentials: %w", err)
	}

	return credentials, nil
}

func (a AuthRepository) Permissions(userId int) (string, error) {
	var roleName string
	err := a.db.QueryRow(`SELECT r.role_name
			FROM users u
			JOIN roles r ON u.role_id = r.role_id
			WHERE u.user_id = $1
		`, userId).Scan(&roleName)

	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("%v with id %d", service.ErrUserNotFound, userId)
		return "", service.ErrUserNotFound
	}

	if err != nil {
		return "", err
	}

	return roleName, nil

}
