package service

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInternalError          = errors.New("internal server error")
	ErrCinemaSessionsNotFound = errors.New("no cinema sessions were found")
	ErrTicketExists           = errors.New("ticket already exists")
)

type Ticket struct {
	Id         int
	MovieName  string
	Date       string
	StartTime  string
	Duration   int
	HallId     int
	SeatNumber int
}

func NewTicketEntity(id, hallId, seat, duration int, movie string, startTime time.Time) Ticket {
	return Ticket{
		Id:         id,
		MovieName:  movie,
		Date:       startTime.Format("2006-01-02"),
		StartTime:  startTime.Format("15:04:05"),
		Duration:   duration,
		HallId:     hallId,
		SeatNumber: seat,
	}
}

type repository interface {
	SessionExists(id int) (bool, error)
	CreateTicket(sessionId, userId, seatNum int) (Ticket, error)
}

type ticketGenerator interface {
	GenerateTicket(t Ticket, outputPath string) error
}

type Service struct {
	r   repository
	gen ticketGenerator
}

func New(r repository, t ticketGenerator) Service {
	return Service{
		r:   r,
		gen: t,
	}
}

func (s Service) BuyTicket(sessionId, userId, seatNum int) (string, error) {
	ok, err := s.r.SessionExists(sessionId)
	if err != nil {
		return "", ErrInternalError
	}

	if !ok {
		return "", ErrCinemaSessionsNotFound
	}

	ticket, err := s.r.CreateTicket(sessionId, userId, seatNum)
	if errors.Is(err, ErrTicketExists) {
		return "", err
	}

	if err != nil {
		return "", ErrInternalError
	}

	outputPath := fmt.Sprintf("ticket%d.pdf", ticket.Id)
	err = s.gen.GenerateTicket(ticket, outputPath)
	if err != nil {
		return "", ErrInternalError
	}
	return outputPath, nil
}