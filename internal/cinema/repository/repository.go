package repository

import "C"
import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type CinemaRepository struct {
	db *sql.DB
}

type Repository interface {
	SessionsForHall(hallId int, timestamp string) ([]CinemaSession, error)
}

func New(db *sql.DB) *CinemaRepository {
	return &CinemaRepository{db: db}
}

type Movie struct {
	ID          int
	Title       string
	Genre       string
	ReleaseDate pq.NullTime
	Duration    int
}

type CinemaSession struct {
	ID        int
	MovieId   int
	StartTime string
	EndTime   string
	Status    string
}

func (c *CinemaRepository) GetMovie(movieID int) (*Movie, error) {
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

func (c *CinemaRepository) SessionsForHall(hallId int, timestamp string) ([]CinemaSession, error) {
	var cinemaSessions []CinemaSession
	rows, err := c.db.Query("SELECT session_id, movie_id, start_time, end_time "+
		"FROM cinema_sessions WHERE hall_id = $1	AND start_time >= $2", hallId, timestamp)
	for rows.Next() {
		var session CinemaSession
		err := rows.Scan(&session.ID, &session.MovieId, &session.StartTime, &session.EndTime)
		session.setStatus(timestamp)
		if err != nil {
			return nil, err
		}
		cinemaSessions = append(cinemaSessions, session)
	}

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no available cinema sessions were found in hall with ID %d", hallId)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cinema session: %w", err)
	}
	return cinemaSessions, nil
}

func (c *CinemaSession) setStatus(timestamp string) {
	if c.StartTime <= timestamp && timestamp <= c.EndTime {
		c.Status = "on_air"
	} else {
		c.Status = "scheduled"
	}
}
