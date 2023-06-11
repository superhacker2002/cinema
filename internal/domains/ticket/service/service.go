package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	ErrInternalError          = errors.New("internal server error")
	ErrCinemaSessionsNotFound = errors.New("no cinema sessions were found")
	ErrTicketExists           = errors.New("ticket already exists")
)

const (
	dateLayout = "2006-01-02"
	timeLayout = "15:04:05"
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
		Date:       startTime.Format(dateLayout),
		StartTime:  startTime.Format(timeLayout),
		Duration:   duration,
		HallId:     hallId,
		SeatNumber: seat,
	}
}

type repository interface {
	SessionExists(id int) (bool, error)
	TicketExists(sessionId, seatNum int) (bool, error)
	CreateTicket(sessionId, userId, seatNum int) (Ticket, error)
}

type ticketGenerator interface {
	GenerateTicket(t Ticket, w io.Writer) error
}

type ticketsStorage interface {
	StoreTicket(file *os.File) (string, error)
}

type Service struct {
	r       repository
	gen     ticketGenerator
	storage ticketsStorage
}

func New(r repository, t ticketGenerator, s ticketsStorage) Service {
	return Service{
		r:       r,
		gen:     t,
		storage: s,
	}
}

func (s Service) BuyTicket(sessionId, userId, seatNum int) (string, error) {
	exists, err := s.r.TicketExists(sessionId, seatNum)
	if err != nil {
		return "", ErrInternalError
	}

	if exists {
		return "", ErrTicketExists
	}

	exists, err = s.r.SessionExists(sessionId)
	if err != nil {
		return "", ErrInternalError
	}

	if !exists {
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
	ticketFile, err := os.Create(outputPath)

	err = s.gen.GenerateTicket(ticket, ticketFile)
	if err != nil {
		return "", ErrInternalError
	}

	ticketFile, err = os.Open(outputPath)
	defer ticketFile.Close()

	path, err := s.storage.StoreTicket(ticketFile)
	if err != nil {
		return "", ErrInternalError
	}

	return path, nil
}
