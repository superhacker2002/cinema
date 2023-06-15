package service

import (
	"errors"
)

var (
	ErrInternalError = errors.New("internal server error")
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotFound  = errors.New("user was not found")
)

const (
	AdminRoleId = 1
	UserRoleId  = 2
)

type repository interface {
	CreateUser(username string, passwordHash string) (userId int, err error)
	MakeAdmin(userId int) (bool, error)
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

func (s Service) MakeAdmin(userId int) error {
	found, err := s.r.MakeAdmin(userId)
	if err != nil {
		return ErrInternalError
	}
	if !found {
		return ErrUserNotFound
	}
	return nil
}
