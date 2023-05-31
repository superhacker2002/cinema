package repository

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasession/entity"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasession/service"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type SessionsRepository struct {
	db *sql.DB
	tz *time.Location
}

func New(db *sql.DB, timeZone *time.Location) *SessionsRepository {
	return &SessionsRepository{db: db, tz: timeZone}
}

type CinemaSession struct {
	ID        int
	MovieId   int
	HallId    int
	StartTime time.Time
	EndTime   time.Time
	Price     float32
}

func (s *SessionsRepository) SessionsForHall(hallId int, date string) ([]entity.CinemaSession, error) {
	rows, err := s.db.Query("SELECT session_id, movie_id, start_time, end_time "+
		"FROM cinema_sessions "+
		"WHERE hall_id = $1 AND date_trunc('day', start_time) = $2 "+
		"ORDER BY start_time ", hallId, date)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get cinema sessions: %w", err)
	}

	cinemaSessions, err := s.readCinemaSessions(rows)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return cinemaSessions, nil
}

func (s *SessionsRepository) AllSessions(date string, offset, limit int) ([]entity.CinemaSession, error) {
	rows, err := s.db.Query("SELECT session_id, movie_id, start_time, end_time "+
		"FROM cinema_sessions "+
		"WHERE start_time >= $1 "+
		"ORDER BY hall_id, start_time "+
		"OFFSET $2 "+
		"LIMIT $3", date, offset, limit)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get cinema sessions: %w", err)
	}

	cinemaSessions, err := s.readCinemaSessions(rows)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return cinemaSessions, nil
}

func (s *SessionsRepository) readCinemaSessions(rows *sql.Rows) ([]entity.CinemaSession, error) {
	var cinemaSessions []entity.CinemaSession
	for rows.Next() {
		var session CinemaSession
		if err := rows.Scan(&session.ID, &session.MovieId, &session.StartTime, &session.EndTime); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("failed to get cinema session: %w", err)
		}
		cinemaSessions = append(cinemaSessions,
			entity.New(session.ID, session.MovieId, session.HallId, session.StartTime, session.EndTime, s.tz))
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error while iterating over cinema sessions: %w", err)
	}

	if len(cinemaSessions) == 0 {
		log.Println(service.ErrCinemaSessionsNotFound)
		return nil, service.ErrCinemaSessionsNotFound
	}

	return cinemaSessions, nil
}
