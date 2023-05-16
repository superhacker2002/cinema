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

type cinemaSession struct {
	id      int
	movieId int
	status  string
}

func (c CinemaRepository) GetMovie(movieID int) (*Movie, error) {
	var movie Movie

	err := c.db.QueryRow("SELECT movie_id, title, genre, release_date, duration FROM movies WHERE movie_id = $1", movieID).
		Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.ReleaseDate, &movie.Duration)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("movie not found: %w", err)
		}
		return nil, fmt.Errorf("could not get movie: %w", err)
	}

	return &movie, nil
}

func (c CinemaRepository) SessionsForHall(hallId int) error {
	var session cinemaSession
	err := c.db.QueryRow("SELECT session_id, movie_id FROM cinema_sessions WHERE hall_id = $1", hallId).
		Scan(&session.id, &session.movieId)
	fmt.Println(session.id, session.movieId)
	if err != nil {
		return err
	}
	return nil
}
