package service

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/entity"
	"errors"
)

type repository interface {
	SessionsForHall(hallId int, date string) ([]entity.CinemaSession, error)
	AllSessions(date string, offset, limit int) ([]entity.CinemaSession, error)
	CreateSession(movieId, hallId int, startTime string, price float32) (int, error)
	DeleteSession(id int) error
}

var (
	ErrInternalError          = errors.New("internal server error")
	ErrCinemaSessionsNotFound = errors.New("no cinema sessions were found")
	ErrHallIsBusy             = errors.New("hall is busy at the time")
	ErrHallNotFound           = errors.New("hall was not found")
	ErrMovieNotFound          = errors.New("movie was not found")
)

type service struct {
	r repository
}

func New(r repository) service {
	return service{r: r}
}

func (s service) AllSessions(date string, offset, limit int) ([]entity.CinemaSession, error) {
	sessions, err := s.r.AllSessions(date, offset, limit)

	return sessions, err
}

func (s service) SessionsForHall(hallId int, date string) ([]entity.CinemaSession, error) {
	sessions, err := s.r.SessionsForHall(hallId, date)

	return sessions, err
}

func (s service) CreateSession(movieId, hallId int, startTime string, price float32) (int, error) {
	sessionId, err := s.r.CreateSession(movieId, hallId, startTime, price)

	return sessionId, err
}

func (s service) DeleteSession(id int) error {
	return s.r.DeleteSession(id)
}
