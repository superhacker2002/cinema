package repository

import "C"
import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

var ErrCinemaSessionsNotFound = errors.New("no available cinema sessions were found")

type SessionsRepository struct {
	db *sql.DB
}

type Repository interface {
	SessionsForHall(hallId int, date string, offset, limit int) ([]CinemaSession, error)
	AllSessions(timestamp string, offset, limit int) ([]CinemaSession, error)
}

func New(db *sql.DB) *SessionsRepository {
	return &SessionsRepository{db: db}
}

type CinemaSession struct {
	ID        int
	MovieId   int
	StartTime string
	EndTime   string
	Status    string
}

func (c *SessionsRepository) SessionsForHall(hallId int, date string, offset, limit int) ([]CinemaSession, error) {
	rows, err := c.db.Query("SELECT session_id, movie_id, start_time, end_time "+
		"FROM cinema_sessions WHERE hall_id = $1 AND date_trunc('day', start_time) = $2 "+
		"ORDER BY start_time OFFSET $3 LIMIT $4", hallId, date, offset, limit)

	if err != nil {
		return nil, fmt.Errorf("failed to get cinema sessions: %w", err)
	}

	cinemaSessions, err := readCinemaSessions(rows, date)
	if err != nil {
		return nil, err
	}

	return cinemaSessions, nil
}

func (c *SessionsRepository) AllSessions(timestamp string, offset, limit int) ([]CinemaSession, error) {
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
		if err := session.setStatus(); err != nil {
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

func (c *CinemaSession) setStatus() error {
	layout := time.RFC3339
	start, err := time.Parse(layout, c.StartTime)
	if err != nil {
		return err
	}
	end, err := time.Parse(layout, c.EndTime)
	if err != nil {
		return err
	}
	current := time.Now()

	if start.Before(current) && end.After(current) {
		c.Status = "on_air"
	} else if end.Before(current) {
		c.Status = "passed"
	} else {
		c.Status = "scheduled"
	}
	return nil
}
