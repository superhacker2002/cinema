package repository

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/user/service"
	"database/sql"
	"errors"
	"fmt"
)

const UserRoleID = 2

type UserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

func (u UserRepository) CreateUser(username string, passwordHash string) (userId int, err error) {
	err = u.db.QueryRow("SELECT user_id FROM users WHERE username = $1", username).
		Scan(&userId)
	if err == nil {
		return 0, fmt.Errorf("%w: %q", service.ErrUserExists, username)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("failed to check if user with username %q exists: %w",
			username, err)
	}

	err = u.db.QueryRow("INSERT INTO users (username, hashed_password, role_id) "+
		"VALUES ($1, $2, $3) RETURNING user_id", username, passwordHash, UserRoleID).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("could not create user: %w", err)
	}

	return userId, nil
}
