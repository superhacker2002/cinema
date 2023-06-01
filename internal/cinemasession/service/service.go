package service

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasession/entity"
	"errors"
	"fmt"
	"log"
	"sort"
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
	HallIsBusy(sessionId, hallId int, startTime, endTime string) (bool, error)
	UpdateSession(id, movieId, hallId int, startTime, endTime string, price float32) error
	AvailableSeats(sessionId int) ([]int, error)
}

var (
	ErrInternalError          = errors.New("internal server error")
	ErrCinemaSessionsNotFound = errors.New("no cinema sessions were found")
	ErrHallIsBusy             = errors.New("hall is busy at the time")
	ErrHallNotFound           = errors.New("hall was not found")
	ErrMovieNotFound          = errors.New("movie was not found")
	ErrNoAvailableSeats       = errors.New("no available seats found for the cinema session")
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

	hallBusy, err := s.r.HallIsBusy(0, hallId, startTime, endTime)
	if err != nil {
		log.Println(err)
		return 0, ErrInternalError
	}
	if hallBusy {
		return 0, fmt.Errorf("%w at the time %s", ErrHallIsBusy, startTime)
	}

	id, err := s.r.CreateSession(movieId, hallId, startTime, endTime, price)
	if err != nil {
		return 0, ErrInternalError
	}

	return id, nil
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

	err = s.r.DeleteSession(id)
	if err != nil {
		return ErrInternalError
	}

	return nil
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

	hallBusy, err := s.r.HallIsBusy(id, hallId, startTime, endTime)
	if err != nil {
		log.Println(err)
		return ErrInternalError
	}
	if hallBusy {
		return fmt.Errorf("%w at the time %s", ErrHallIsBusy, startTime)
	}

	err = s.r.UpdateSession(id, movieId, hallId, startTime, endTime, price)
	if err != nil {
		return ErrInternalError
	}

	return nil
}

func (s Service) AvailableSeats(sessionId int) ([]int, error) {
	ok, err := s.r.SessionExists(sessionId)
	if err != nil {
		return nil, ErrInternalError
	}
	if !ok {
		return nil, ErrCinemaSessionsNotFound
	}
	seats, err := s.r.AvailableSeats(sessionId)
	if errors.Is(err, ErrNoAvailableSeats) {
		return nil, err
	}
	if err != nil {
		return nil, ErrInternalError
	}
	sort.Ints(seats)
	return seats, nil
}
