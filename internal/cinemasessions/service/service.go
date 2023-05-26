package service

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/entity"
	"errors"
	"fmt"
	"log"
)

type repository interface {
	SessionsForHall(hallId int, date string) ([]entity.CinemaSession, error)
	AllSessions(date string, offset, limit int) ([]entity.CinemaSession, error)
	CreateSession(movieId, hallId int, startTime, endTime string, price float32) (int, error)
	DeleteSession(id int) error
	SessionExists(id int) (bool, error)
	HallExists(id int) (bool, error)
	MovieExists(id int) (bool, error)
	SessionEndTime(id int, startTime string) (string, error)
	HallIsBusy(movieId, hallId int, startTime, endTime string) (bool, error)
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

func (s service) CreateSession(movieId, hallId int, startTime string, price float32) (int, error) {
	hallExists, err := s.r.HallExists(hallId)
	if err != nil {
		log.Println(err)
		return 0, ErrInternalError
	}
	if !hallExists {
		return 0, ErrHallNotFound
	}

	movieExists, err := s.r.MovieExists(hallId)
	if err != nil {
		log.Println(err)
		return 0, ErrInternalError
	}
	if !movieExists {
		return 0, ErrMovieNotFound
	}

	endTime, err := s.r.SessionEndTime(movieId, startTime)
	if err != nil {
		log.Println(err)
		return 0, ErrInternalError
	}

	hallBusy, err := s.r.HallIsBusy(movieId, hallId, startTime, endTime)
	if err != nil {
		log.Println(err)
		return 0, ErrInternalError
	}
	if hallBusy {
		return 0, fmt.Errorf("%w at the time %s", ErrHallIsBusy, startTime)
	}

	return s.r.CreateSession(movieId, hallId, startTime, endTime, price)
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
