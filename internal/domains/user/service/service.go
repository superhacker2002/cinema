package service

import (
	"errors"
)

var (
	ErrInternalError = errors.New("internal server error")
	ErrUserExists    = errors.New("user already exists")
)

type repository interface {
	CreateUser(username string, passwordHash string) (userId int, err error)
}

type Service struct {
	r repository
}

func New(r repository) Service {
	return Service{r: r}
}

func (s Service) CreateUser(username string, passwordHash string) (userId int, err error) {
	id, err := s.r.CreateUser(username, passwordHash)
	if errors.Is(err, ErrUserExists) {
		return 0, err
	}

	if err != nil {
		return 0, ErrInternalError
	}

	return id, nil
}
