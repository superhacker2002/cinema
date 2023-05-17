package repository

import "C"
import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"time"
)

var ErrCinemaSessionsNotFound = errors.New("no available cinema sessions were found")

type CinemaRepository struct {
	db *sql.DB
}

type Repository interface {
	SessionsForHall(hallId int, timestamp string, offset, limit int) ([]CinemaSession, error)
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

func (c *CinemaRepository) SessionsForHall(hallId int, timestamp string, offset, limit int) ([]CinemaSession, error) {
	log.Println(timestamp)
	rows, err := c.db.Query("SELECT session_id, movie_id, start_time, end_time "+
		"FROM cinema_sessions WHERE hall_id = $1 AND end_time > $2"+
		"ORDER BY start_time OFFSET $3 LIMIT $4", hallId, timestamp, offset, limit)

	if err != nil {
		return nil, fmt.Errorf("failed to get cinema sessions: %w", err)
	}

	cinemaSessions, err := readCinemaSessions(rows, timestamp)
	if err != nil {
		return nil, err
	}

	return cinemaSessions, nil
}

func (c *CinemaRepository) AllSessions(timestamp string, offset, limit int) ([]CinemaSession, error) {
	log.Println(timestamp)
	rows, err := c.db.Query("SELECT session_id, movie_id, start_time, end_time "+
		"FROM cinema_sessions WHERE end_time > $1 ORDER BY hall_id, start_time OFFSET $2 LIMIT $3",
		timestamp, offset, limit)

	if err != nil {
		return nil, fmt.Errorf("failed to get cinema sessions: %w", err)
	}

	cinemaSessions, err := readCinemaSessions(rows, timestamp)
	if err != nil {
		return nil, err
	}

	return cinemaSessions, nil
}

func readCinemaSessions(rows *sql.Rows, timestamp string) ([]CinemaSession, error) {
	var cinemaSessions []CinemaSession
	for rows.Next() {
		var session CinemaSession
		if err := rows.Scan(&session.ID, &session.MovieId, &session.StartTime, &session.EndTime); err != nil {
			return nil, fmt.Errorf("failed to get cinema session: %w", err)
		}
		if err := session.setStatus(timestamp); err != nil {
			return nil, fmt.Errorf("failed to set cinema session status: %w", err)
		}
		cinemaSessions = append(cinemaSessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over cinema sessions: %w", err)
	}

	if len(cinemaSessions) == 0 {
		return nil, ErrCinemaSessionsNotFound
	}

	return cinemaSessions, nil
}

func (c *CinemaSession) setStatus(timestamp string) error {
	start, err := time.Parse(time.RFC3339, c.StartTime)
	if err != nil {
		return err
	}
	end, err := time.Parse(time.RFC3339, c.EndTime)
	if err != nil {
		return err
	}
	current, err := time.Parse("2006-01-02 15:04:05", timestamp)
	if err != nil {
		return err
	}

	if start.Before(current) && end.After(current) {
		c.Status = "on_air"
	} else {
		c.Status = "scheduled"
	}
	return nil
}
