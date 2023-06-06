package service

import (
	"errors"
)

type auth interface {
	Authenticate(username string, passwordHash string) (token string, err error)
	VerifyToken(token string) (userID int, err error)
}

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
