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
