package repository

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasession/entity"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasession/service"
	"database/sql"
	"fmt"
	"log"
	"strconv"
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
	rows, err := s.db.Query(`SELECT session_id, movie_id, hall_id, start_time, end_time, price
		FROM cinema_sessions
		WHERE hall_id = $1 AND date_trunc('day', start_time) = $2
		ORDER BY start_time`, hallId, date)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get cinema sessions: %w", err)
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	cinemaSessions, err := s.readCinemaSessions(rows)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return cinemaSessions, nil
}

func (s *SessionsRepository) AllSessions(date string, offset, limit int) ([]entity.CinemaSession, error) {
	rows, err := s.db.Query(`SELECT session_id, movie_id, hall_id, start_time, end_time, price
		FROM cinema_sessions
		WHERE start_time >= $1
		ORDER BY hall_id, start_time
		OFFSET $2
		LIMIT $3`, date, offset, limit)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get cinema sessions: %w", err)
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	cinemaSessions, err := s.readCinemaSessions(rows)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return cinemaSessions, nil
}

func (s *SessionsRepository) CreateSession(movieId, hallId int, startTime, endTime string, price float32) (sessionId int, err error) {
	err = s.db.QueryRow(`INSERT INTO cinema_sessions (movie_id, hall_id, start_time, end_time, price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING session_id`, movieId, hallId, startTime, endTime, price).Scan(&sessionId)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return sessionId, nil
}

func (s *SessionsRepository) HallIsBusy(sessionId, hallId int, startTime, endTime string) (bool, error) {
	row := s.db.QueryRow(`SELECT session_id
		FROM cinema_sessions
		WHERE hall_id = $1 AND session_id != $2 AND (start_time BETWEEN $3 AND $4
		OR (start_time <= $3 AND end_time > $3))`, hallId, sessionId, startTime, endTime)

	var sessionExistId int
	err := row.Scan(&sessionExistId)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("failed to check if cinema session can be created: %w", err)
	}

	return err == nil, nil
}

func (s *SessionsRepository) SessionEndTime(id int, startTime string) (string, error) {
	var (
		duration string
		endTime  string
	)
	row := s.db.QueryRow("SELECT duration FROM movies WHERE movie_id = $1", id)
	if err := row.Scan(&duration); err == sql.ErrNoRows {
		return endTime, fmt.Errorf("%w with id %d", service.ErrMovieNotFound, id)
	} else if err != nil {
		return endTime, fmt.Errorf("failed to get session end time: %w", err)
	}

	const layout = "2006-01-02 15:04:05 MST"
	start, err := time.Parse(layout, startTime)
	if err != nil {
		return endTime, fmt.Errorf("failed to get session end time: %w", err)
	}

	durationMinutes, err := strconv.Atoi(duration)
	if err != nil {
		return endTime, fmt.Errorf("failed to get session end time: %w", err)
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

func (s *SessionsRepository) UpdateSession(id, movieId, hallId int, startTime, endTime string, price float32) error {
	_, err := s.db.Exec(`UPDATE cinema_sessions
		SET movie_id = $1, hall_id = $2, start_time = $3, end_time = $4, price = $5
		WHERE session_id = $6`, movieId, hallId, startTime, endTime, price, id)

	if err != nil {
		return fmt.Errorf("failed to update cinema session: %w", err)
	}

	return nil
}

func (s *SessionsRepository) AvailableSeats(sessionId int) ([]int, error) {
	rows, err := s.db.Query(`SELECT seat_number
				FROM (
					SELECT generate_series(1, capacity) AS seat_number
					FROM halls
					WHERE hall_id = (SELECT hall_id FROM cinema_sessions WHERE session_id = $1)
				) AS all_seats
				EXCEPT (
					SELECT seat_number
					FROM tickets
					WHERE session_id = $1
				)`, sessionId)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get available seats: %w", err)
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	var freeSeats []int
	for rows.Next() {
		var seatNumber int
		err := rows.Scan(&seatNumber)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("failed to get available seats: %w", err)
		}
		freeSeats = append(freeSeats, seatNumber)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error while iterating over halls seats: %w", err)
	}

	if len(freeSeats) == 0 {
		log.Println(service.ErrNoAvailableSeats)
		return nil, service.ErrNoAvailableSeats
	}

	return freeSeats, nil
}

func (s *SessionsRepository) SessionExists(id int) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM cinema_sessions WHERE session_id = $1", id).Scan(&count)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("failed to check if session exists %w", err)
	}

	return count > 0, nil
}

func (s *SessionsRepository) HallExists(id int) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM halls WHERE hall_id = $1", id).Scan(&count)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("failed to check if hall exists %w", err)
	}

	return count > 0, nil
}

func (s *SessionsRepository) MovieExists(id int) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM movies WHERE movie_id = $1", id).Scan(&count)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("failed to check if movie exists %w", err)
	}

	return count > 0, nil
}

func (s *SessionsRepository) readCinemaSessions(rows *sql.Rows) ([]entity.CinemaSession, error) {
	var cinemaSessions []entity.CinemaSession
	for rows.Next() {
		var session CinemaSession
		if err := rows.Scan(&session.ID, &session.MovieId, &session.HallId,
			&session.StartTime, &session.EndTime, &session.Price); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("failed to get cinema session: %w", err)
		}
		cinemaSessions = append(cinemaSessions,
			entity.New(session.ID, session.MovieId, session.HallId, session.StartTime, session.EndTime, session.Price, s.tz))
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
