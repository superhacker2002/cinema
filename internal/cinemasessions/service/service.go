package service

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/entity"
	"errors"
	"log"
)

type repository interface {
	SessionsForHall(hallId int, date string) ([]entity.CinemaSession, error)
	AllSessions(date string, offset, limit int) ([]entity.CinemaSession, error)
	CreateSession(movieId, hallId int, startTime string, price float32) (int, error)
	DeleteSession(id int) error
	SessionExists(id int) (bool, error)
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
	return s.r.AllSessions(date, offset, limit)
}

func (s service) SessionsForHall(hallId int, date string) ([]entity.CinemaSession, error) {
	return s.r.SessionsForHall(hallId, date)
}

func (s service) CreateSession(movieId, hallId int, startTime string, price float32) (int, error) {
	return s.r.CreateSession(movieId, hallId, startTime, price)
}

func (s service) DeleteSession(id int) error {
	exists, err := s.r.SessionExists(id)
	if err != nil {
		log.Println(err)
		return ErrInternalError
	}

	if !exists {
		return ErrCinemaSessionsNotFound
	}

	return s.r.DeleteSession(id)
}
