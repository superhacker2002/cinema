package repository

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/entity"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/service"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

var timeZone = time.FixedZone("UTC+4", 4*60*60)

type SessionsRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *SessionsRepository {
	return &SessionsRepository{db: db}
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

	cinemaSessions, err := readCinemaSessions(rows)
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

	cinemaSessions, err := readCinemaSessions(rows)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return cinemaSessions, nil
}

func (s *SessionsRepository) CreateSession(movieId, hallId int, startTime string, price float32) (sessionId int, err error) {
	exists, err := s.hallExists(hallId)
	if err != nil {
		log.Println(err)
		return 0, fmt.Errorf("failed to check if hall with id %d exists: %w", hallId, err)
	}
	if !exists {
		return 0, fmt.Errorf("%w: id %d", service.ErrHallNotFound, hallId)
	}

	endTime, err := s.sessionEndTime(movieId, startTime)
	if err != nil {
		log.Println(err)
		return 0, fmt.Errorf("failed to get session end time: %w", err)
	}

	hallBusy, err := s.hallIsBusy(movieId, hallId, startTime, endTime)
	if err != nil {
		log.Println(err)
		return 0, fmt.Errorf("failed to check if cinema session can be created: %w", err)
	}
	if hallBusy {
		return 0, fmt.Errorf("%w at the time %s", service.ErrHallIsBusy, startTime)
	}

	log.Println(startTime)
	log.Println(endTime)
	err = s.db.QueryRow("INSERT INTO cinema_sessions (movie_id, hall_id, start_time, end_time, price)"+
		"VALUES ($1, $2, $3, $4, $5)"+
		"RETURNING session_id", movieId, hallId, startTime, endTime, price).Scan(&sessionId)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	var start string
	_ = s.db.QueryRow("SELECT start_time FROM cinema_sessions WHERE session_id = $1", sessionId).Scan(&start)
	log.Println(start)

	return sessionId, nil
}

func (s *SessionsRepository) hallIsBusy(movieId, hallId int, startTime, endTime string) (bool, error) {
	row := s.db.QueryRow("SELECT session_id "+
		"FROM cinema_sessions "+
		"WHERE hall_id = $1 AND $2 < start_time AND start_time < $3 "+
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
		return endTime, fmt.Errorf("%w with id %d", service.ErrMovieNotFound, id)
	} else if err != nil {
		return endTime, err
	}

	const layout = "2006-01-02 15:04:05 MST"
	start, err := time.Parse(layout, startTime)
	if err != nil {
		return endTime, err
	}

	durationMinutes, err := strconv.Atoi(duration)
	if err != nil {
		return endTime, err
	}
	endTime = start.Add(time.Minute * time.Duration(durationMinutes)).Format(layout)
	return endTime, nil
}

func (s *SessionsRepository) DeleteSession(id int) error {
	_, err := s.db.Exec("DELETE FROM cinema_sessions WHERE session_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete cinema session: %w", err)
	}

	return nil
}

func (s *SessionsRepository) SessionExists(id int) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM cinema_sessions WHERE session_id = $1", id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if session exists %w", err)
	}

	return count > 0, nil
}

func (s *SessionsRepository) hallExists(id int) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM halls WHERE hall_id = $1", id).Scan(&count)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return count > 0, nil
}

func readCinemaSessions(rows *sql.Rows) ([]entity.CinemaSession, error) {
	var cinemaSessions []entity.CinemaSession
	for rows.Next() {
		var session CinemaSession
		if err := rows.Scan(&session.ID, &session.MovieId, &session.StartTime, &session.EndTime); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("failed to get cinema session: %w", err)
		}
		session.StartTime = session.StartTime.In(timeZone)
		session.EndTime = session.EndTime.In(timeZone)
		cinemaSessions = append(cinemaSessions,
			entity.New(session.ID, session.MovieId, session.StartTime, session.EndTime))
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
