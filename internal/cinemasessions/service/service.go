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
	HallIsBusy(movieId, hallId int, startTime, endTime string, sessionId int) (bool, error)
	UpdateSession(id, movieId, hallId int, startTime, endTime string, price float32) error
}

var (
	ErrInternalError          = errors.New("internal server error")
	ErrCinemaSessionsNotFound = errors.New("no cinema sessions were found")
	ErrHallIsBusy             = errors.New("hall is busy at the time")
	ErrHallNotFound           = errors.New("hall was not found")
	ErrMovieNotFound          = errors.New("movie was not found")
)

type Service struct {
	r repository
}

func New(r repository) Service {
	return Service{r: r}
}

func (s Service) AllSessions(date string, offset, limit int) ([]entity.CinemaSession, error) {
	sessions, err := s.r.AllSessions(date, offset, limit)
	if errors.Is(err, ErrCinemaSessionsNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, ErrInternalError
	}
	return sessions, nil
}

func (s Service) SessionsForHall(hallId int, date string) ([]entity.CinemaSession, error) {
	sessions, err := s.r.SessionsForHall(hallId, date)
	if errors.Is(err, ErrCinemaSessionsNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, ErrInternalError
	}
	return sessions, nil
}

func (s Service) CreateSession(movieId, hallId int, startTime string, price float32) (int, error) {
	hallExists, err := s.r.HallExists(hallId)
	if err != nil {
		log.Println(err)
		return 0, ErrInternalError
	}
	if !hallExists {
		return 0, ErrHallNotFound
	}

	movieExists, err := s.r.MovieExists(movieId)
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

	hallBusy, err := s.r.HallIsBusy(movieId, hallId, startTime, endTime, 0)
	if err != nil {
		log.Println(err)
		return 0, ErrInternalError
	}
	if hallBusy {
		return 0, fmt.Errorf("%w at the time %s", ErrHallIsBusy, startTime)
	}

	return s.r.CreateSession(movieId, hallId, startTime, endTime, price)
}

func (s Service) DeleteSession(id int) error {
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

func (s Service) UpdateSession(id, movieId, hallId int, startTime string, price float32) error {
	ok, err := s.r.SessionExists(id)
	if err != nil {
		log.Println(err)
		return ErrInternalError
	}

	if !ok {
		return ErrCinemaSessionsNotFound
	}

	ok, err = s.r.HallExists(hallId)
	if err != nil {
		log.Println(err)
		return ErrInternalError
	}
	if !ok {
		return ErrHallNotFound
	}

	ok, err = s.r.MovieExists(movieId)
	if err != nil {
		log.Println(err)
		return ErrInternalError
	}
	if !ok {
		return ErrMovieNotFound
	}

	endTime, err := s.r.SessionEndTime(movieId, startTime)
	if err != nil {
		log.Println(err)
		return ErrInternalError
	}

	hallBusy, err := s.r.HallIsBusy(movieId, hallId, startTime, endTime, id)
	if err != nil {
		log.Println(err)
		return ErrInternalError
	}
	if hallBusy {
		return fmt.Errorf("%w at the time %s", ErrHallIsBusy, startTime)
	}

	return s.r.UpdateSession(id, movieId, hallId, startTime, endTime, price)
}
