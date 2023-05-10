package repository

import (
	"database/sql"
	"fmt"

	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinema/handler"
)

type repository struct {
	db *sql.DB
}

type Movie = handler.Movie // Import the Movie struct from the handler package

func New(db *sql.DB) repository {
	return repository{db: db}
}

func NewMovieRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetMovie(movieID int) (*Movie, error) {
	var movie Movie

	err := r.db.QueryRow("SELECT id, title, genre, release_date, duration FROM movies WHERE id = $1", movieID).
		Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.ReleaseDate, &movie.Duration)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("movie not found")
		}
		return nil, err
	}
	return &movie, nil
}
