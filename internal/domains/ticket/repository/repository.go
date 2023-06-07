package repository

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/ticket/service"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type ticket struct {
	Id        int
	MovieName string
	StartTime time.Time
	Duration  int
	HallId    int
}

type TicketRepository struct {
	db *sql.DB
}

func New(db *sql.DB) TicketRepository {
	return TicketRepository{db: db}
}

func (t TicketRepository) CreateTicket(sessionId, userId, seatNum int) (service.Ticket, error) {
	sessionTicket, err := t.sessionInfo(sessionId)
	if err != nil {
		return service.Ticket{}, err
	}

	exists, err := t.ticketExists(sessionId, seatNum)
	if err != nil {
		log.Println(err)
		return service.Ticket{}, err
	}

	if exists {
		log.Printf("%v for the session with id %d and seat number %d", service.ErrTicketExists, sessionId, seatNum)
		return service.Ticket{}, service.ErrTicketExists
	}

	var id int
	err = t.db.QueryRow(`INSERT INTO tickets (session_id, user_id, seat_number) 
				VALUES ($1, $2, $3) RETURNING ticket_id`, sessionId, userId, seatNum).Scan(&id)

	if err != nil {
		log.Println(err)
		return service.Ticket{}, err
	}

	return service.NewTicketEntity(id, sessionTicket.HallId, seatNum,
		sessionTicket.Duration, sessionTicket.MovieName, sessionTicket.StartTime), nil
}

func (t TicketRepository) SessionExists(id int) (bool, error) {
	var count int
	err := t.db.QueryRow("SELECT COUNT(*) FROM cinema_sessions WHERE session_id = $1", id).Scan(&count)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("failed to check if session exists %w", err)
	}

	return count > 0, nil
}

func (t TicketRepository) sessionInfo(sessionId int) (ticket, error) {
	var sessionTicket ticket
	err := t.db.QueryRow(`
		SELECT m.title, s.start_time, m.duration, s.hall_id
		FROM cinema_sessions s
		JOIN movies m ON s.movie_id = m.movie_id
		WHERE s.session_id = $1`, sessionId).Scan(&sessionTicket.MovieName,
		&sessionTicket.StartTime, &sessionTicket.Duration, &sessionTicket.HallId)
	if err != nil {
		log.Println(err)
		return ticket{}, err
	}
	return sessionTicket, nil
}

func (t TicketRepository) ticketExists(sessionId, seatNum int) (bool, error) {
	var count int
	err := t.db.QueryRow(`SELECT COUNT(*) 
				FROM tickets 
				WHERE session_id = $1 AND seat_number = $2`, sessionId, seatNum).Scan(&count)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("failed to check if ticket exists %w", err)
	}

	return count > 0, nil
}
