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
	var roleId int
	err := a.db.QueryRow("SELECT role_id FROM users WHERE user_id = $1", userId).Scan(&roleId)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("could not get user permissions:", err)
		return "", service.ErrUserNotFound
	}

	if err != nil {
		return "", err
	}

	var roleName string
	err = a.db.QueryRow("SELECT role_name FROM roles WHERE role_id = $1", roleId).Scan(&roleName)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("could not get user role:", err)
		return "", err
	}

	return roleName, nil
}
