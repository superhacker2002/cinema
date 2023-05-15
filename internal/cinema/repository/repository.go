package repository

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type CinemaRepository struct {
	db *sql.DB
}

func New(db *sql.DB) CinemaRepository {
	return CinemaRepository{db: db}
}

type Movie struct {
	ID          int
	Title       string
	Genre       string
	ReleaseDate pq.NullTime
	Duration    int
}

func (r CinemaRepository) GetMovie(movieID int) (*Movie, error) {
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
