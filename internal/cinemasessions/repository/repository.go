package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	ErrCinemaSessionsNotFound = errors.New("no cinema sessions were found")
	ErrHallIsBusy             = errors.New("hall is busy")
)

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
	CreateSession(movieId, hallId int, startTime string, price float32) (sessionId int, err error)
	DeleteSession(id int) error
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

func (s *SessionsRepository) CreateSession(movieId, hallId int, startTime string, price float32) (sessionId int, err error) {
	endTime, err := s.sessionEndTime(movieId, startTime)
	if err != nil {
		return 0, fmt.Errorf("failed to get session end time: %w", err)
	}

	hallBusy, err := s.hallIsBusy(movieId, hallId, startTime, endTime)
	if err != nil {
		return 0, fmt.Errorf("failed to check if cinema session can be created: %w", err)
	}
	if hallBusy {
		return 0, fmt.Errorf("%w at the time %s", ErrHallIsBusy, startTime)
	}

	err = s.db.QueryRow("INSERT INTO cinema_sessions (movie_id, hall_id, start_time, end_time) VALUES ($1, $2, $3, $4)"+
		"RETURNING session_id", movieId, hallId, startTime, price).Scan(&sessionId)
	if err != nil {
		return 0, err
	}

	return sessionId, nil
}

func (s *SessionsRepository) hallIsBusy(movieId, hallId int, startTime, endTime string) (bool, error) {
	row := s.db.QueryRow("SELECT session_id"+
		"FROM cinema_sessions"+
		"WHERE hall_id = $1 AND $2 < start_time < $3 "+
		"OR (start_time <= $2 AND end_time > $2)",
		hallId, startTime, endTime)

	var sessionId int
	if err := row.Scan(&sessionId); err == nil {
		return true, nil
	} else if err != sql.ErrNoRows {
		return true, err
	}
	return false, nil
}

func (s *SessionsRepository) sessionEndTime(id int, startTime string) (string, error) {
	var (
		duration string
		endTime  string
	)
	row := s.db.QueryRow("SELECT duration FROM movies WHERE movie_id = $1", id)
	if err := row.Scan(&duration); err == sql.ErrNoRows {
		return endTime, errors.New("movie was not found")
	} else if err != nil {
		return endTime, err
	}

	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return endTime, err
	}

	durationMinutes, err := strconv.Atoi(duration)
	if err != nil {
		return endTime, err
	}
	endTime = start.Add(time.Minute * time.Duration(durationMinutes)).Format("2006-01-02")
	return endTime, nil
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
