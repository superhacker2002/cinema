package service

import (
	"errors"
)

var (
	ErrInternalError          = errors.New("internal server error")
	ErrCinemaSessionsNotFound = errors.New("no cinema sessions were found")
	ErrTicketExists           = errors.New("ticket already exists")
)

type Ticket struct {
	Id         int
	MovieName  string
	StartTime  string
	Duration   int
	HallId     int
	SeatNumber int
}

func NewTicketEntity(id, hallId, seat, duration int, movie, startTime string) Ticket {
	return Ticket{
		Id:         id,
		MovieName:  movie,
		StartTime:  startTime,
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
	GenerateTicket(t Ticket, outputPath string) (string, error)
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

	ticketPath, err := s.gen.GenerateTicket(ticket, "ticket.pdf")
	if err != nil {
		return "", ErrInternalError
	}
	return ticketPath, nil
}
