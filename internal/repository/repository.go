package repository

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/auth"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return Repository{db: db}
}

func NewMovieRepository(db *sql.DB) Repository {
	return Repository{db: db}
}

type Movie struct {
	ID          int
	Title       string
	Genre       string
	ReleaseDate pq.NullTime
	Duration    int
}

func (r Repository) User(username string) (auth.Credentials, error) {
	credentials := auth.Credentials{}
	err := r.db.QueryRow("SELECT user_id, hashed_password FROM users WHERE username=$1", username).
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

func (r Repository) CreateUser(username string, passwordHash string, role string) (userId string, err error) {
	err = r.db.QueryRow("SELECT user_id FROM users WHERE username = $1", username).
		Scan(&userId)
	if err == nil {
		return "", fmt.Errorf("user with username %q already exists", username)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("failed to check if user with username %q exists: %w",
			username, err)
	}

	var roleID string
	err = r.db.QueryRow("SELECT role_id FROM roles WHERE role_name = $1", role).
		Scan(&roleID)
	if err != nil {
		return "", fmt.Errorf("could not get role id: %w", err)
	}

	err = r.db.QueryRow("INSERT INTO users (username, hashed_password, role_id) "+
		"VALUES ($1, $2, $3) RETURNING user_id", username, passwordHash, roleID).Scan(&userId)
	if err != nil {
		return "", fmt.Errorf("could not create user: %w", err)
	}

	return userId, nil
}

func (r Repository) GetMovie(movieID int) (*Movie, error) {
	var movie Movie

	err := r.db.QueryRow("SELECT id, title, genre, release_date, duration FROM movies WHERE id = $1", movieID).
		Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.ReleaseDate, &movie.Duration)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("movie not found: %w", err)
		}
		return nil, fmt.Errorf("could not get movie: %w", err)
	}

	return &movie, nil
}
