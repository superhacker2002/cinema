package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrCinemaSessionsNotFound = errors.New("no cinema sessions were found")

const (
	StatusPassed    = "passed"
	StatusOnAir     = "on_air"
	StatusScheduled = "scheduled"
)

type SessionsRepository struct {
	db *sql.DB
}

type Repository interface {
	SessionsForHall(hallId int, date string) ([]CinemaSession, error)
	AllSessions(date string, offset, limit int) ([]CinemaSession, error)
	DeleteSession(id int) error
}

func New(db *sql.DB) *SessionsRepository {
	return &SessionsRepository{db: db}
}

type CinemaSession struct {
	ID        int
	MovieId   int
	StartTime time.Time
	EndTime   time.Time
	Status    string
}

func (s *SessionsRepository) SessionsForHall(hallId int, date string) ([]CinemaSession, error) {
	rows, err := s.db.Query("SELECT session_id, movie_id, start_time, end_time "+
		"FROM cinema_sessions "+
		"WHERE hall_id = $1 AND date_trunc('day', start_time) = $2 "+
		"ORDER BY start_time ", hallId, date)

	if err != nil {
		return nil, fmt.Errorf("failed to get cinema sessions: %w", err)
	}

	cinemaSessions, err := readCinemaSessions(rows)
	if err != nil {
		return nil, err
	}

	return cinemaSessions, nil
}

func (s *SessionsRepository) AllSessions(date string, offset, limit int) ([]CinemaSession, error) {
	rows, err := s.db.Query("SELECT session_id, movie_id, start_time, end_time "+
		"FROM cinema_sessions "+
		"WHERE start_time >= $1 "+
		"ORDER BY hall_id, start_time "+
		"OFFSET $2 "+
		"LIMIT $3", date, offset, limit)

	if err != nil {
		return nil, fmt.Errorf("failed to get cinema sessions: %w", err)
	}

	cinemaSessions, err := readCinemaSessions(rows)
	if err != nil {
		return nil, err
	}

	return cinemaSessions, nil
}

func (s *SessionsRepository) DeleteSession(id int) error {
	exists, err := s.sessionExists(id)
	if err != nil {
		return fmt.Errorf("failed to check if cinema session with id %d exists: %w", id, err)
	}

	if !exists {
		return fmt.Errorf("%w with id %d", ErrCinemaSessionsNotFound, id)
	}

	_, err = s.db.Exec("DELETE FROM cinema_sessions WHERE session_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete cinema session: %w", err)
	}

	return nil
}

func (s *SessionsRepository) sessionExists(id int) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM cinema_sessions WHERE session_id = $1", id).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func readCinemaSessions(rows *sql.Rows) ([]CinemaSession, error) {
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
	current := time.Now().UTC()

	if c.StartTime.Before(current) && c.EndTime.After(current) {
		c.Status = StatusOnAir
	} else if c.EndTime.Before(current) {
		c.Status = StatusPassed
	} else {
		c.Status = StatusScheduled
	}
	return nil
}
