package service

import "errors"

var (
	ErrHallNotFound  = errors.New("hall not found")
	ErrInternalError = errors.New("internal server error")
)

type Hall struct {
	Id       int
	Name     string
	Capacity int
}

func NewHallEntity(id int, name string, capacity int) Hall {
	return Hall{
		Id:       id,
		Name:     name,
		Capacity: capacity,
	}
}

type repository interface {
	Halls() ([]Hall, error)
	HallById(id int) (Hall, error)
	CreateHall(name string, capacity int) (hallId int, err error)
	UpdateHall(id int, name string, capacity int) error
	DeleteHall(id int) error
	HallExists(id int) (bool, error)
}

type Service struct {
	R repository
}

func (s Service) Halls() ([]Hall, error) {
	return s.R.Halls()
}

func (s Service) HallById(id int) (Hall, error) {
	hall, err := s.R.HallById(id)
	if errors.Is(err, ErrHallNotFound) {
		return Hall{}, ErrHallNotFound
	}
	if err != nil {
		return Hall{}, ErrInternalError
	}
	return hall, nil
}

func (s Service) CreateHall(name string, capacity int) (hallId int, err error) {
	return s.R.CreateHall(name, capacity)
}

func (s Service) UpdateHall(id int, name string, capacity int) error {
	ok, err := s.R.HallExists(id)
	if err != nil {
		return ErrInternalError
	}
	if !ok {
		return ErrHallNotFound
	}
	return s.R.UpdateHall(id, name, capacity)
}

func (s Service) DeleteHall(id int) error {
	ok, err := s.R.HallExists(id)
	if err != nil {
		return ErrInternalError
	}
	if !ok {
		return ErrHallNotFound
	}
	return s.R.DeleteHall(id)
}
