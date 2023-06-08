package service

import "errors"

var (
	ErrHallNotFound  = errors.New("hall not found")
	ErrInternalError = errors.New("internal server error")
)

const AdminRole = "admin"

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
	UpdateHall(id int, name string, capacity int) (found bool, err error)
	DeleteHall(id int) (bool, error)
}

type Service struct {
	r repository
}

func New(r repository) Service {
	return Service{r: r}
}

func (s Service) Halls() ([]Hall, error) {
	halls, err := s.r.Halls()
	if err != nil {
		return []Hall{}, ErrInternalError
	}
	return halls, nil
}

func (s Service) HallById(id int) (Hall, error) {
	hall, err := s.r.HallById(id)
	if errors.Is(err, ErrHallNotFound) {
		return Hall{}, ErrHallNotFound
	}
	if err != nil {
		return Hall{}, ErrInternalError
	}
	return hall, nil
}

func (s Service) CreateHall(name string, capacity int) (hallId int, err error) {
	id, err := s.r.CreateHall(name, capacity)
	if err != nil {
		return 0, ErrInternalError
	}
	return id, nil
}

func (s Service) UpdateHall(id int, name string, capacity int) error {
	found, err := s.r.UpdateHall(id, name, capacity)
	if err != nil {
		return ErrInternalError
	}
	if !found {
		return ErrHallNotFound
	}
	return nil
}

func (s Service) DeleteHall(id int) error {
	found, err := s.r.DeleteHall(id)
	if err != nil {
		return ErrInternalError
	}
	if !found {
		return ErrHallNotFound
	}
	return nil
}
