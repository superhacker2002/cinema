package service

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/entity"
	"errors"
)

type repository interface {
	SessionsForHall(hallId int, date string) ([]entity.CinemaSession, error)
	AllSessions(date string, offset, limit int) ([]entity.CinemaSession, error)
}

var (
	ErrInternalError          = errors.New("internal server error")
	ErrCinemaSessionsNotFound = errors.New("no cinema sessions were found")
)

type service struct {
	r repository
}

func New(r repository) service {
	return service{r: r}
}

func (s service) AllSessions(date string, offset, limit int) ([]entity.CinemaSession, error) {
	sessions, err := s.r.AllSessions(date, offset, limit)
	if errors.Is(err, ErrCinemaSessionsNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, ErrInternalError
	}
	return sessions, nil
}

func (s service) SessionsForHall(hallId int, date string) ([]entity.CinemaSession, error) {
	sessions, err := s.r.SessionsForHall(hallId, date)
	if errors.Is(err, ErrCinemaSessionsNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, ErrInternalError
	}
	return sessions, nil
}
