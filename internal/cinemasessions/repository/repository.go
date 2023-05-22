package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrCinemaSessionsNotFound = errors.New("no available cinema sessions were found")

type Status int

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

func (s *SessionsRepository) deleteSession(id int) error {

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

func (s *CinemaSession) setStatus() error {
	layout := time.RFC3339
	start, err := time.Parse(layout, s.StartTime)
	if err != nil {
		return err
	}
	end, err := time.Parse(layout, s.EndTime)
	if err != nil {
		return err
	}
	current := time.Now().UTC()

	if start.Before(current) && end.After(current) {
		s.Status = StatusOnAir
	} else if end.Before(current) {
		s.Status = StatusPassed
	} else {
		s.Status = StatusScheduled
	}
	return nil
}
