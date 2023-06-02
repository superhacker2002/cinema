package service

import (
	"errors"
	"time"
)

var (
	ErrMovieNotFound = errors.New("movie not found")
	ErrInternalError = errors.New("internal server error")
)

type Movie struct {
	Id          int
	Title       string
	Genre       string
	ReleaseDate string
	Duration    int
}

func NewMovieEntity(id int, title, genre string, releaseDate string, duration int) Movie {
	return Movie{
		Id:          id,
		Title:       title,
		Genre:       genre,
		ReleaseDate: releaseDate,
		Duration:    duration,
	}
}

type repository interface {
	Movies(date string) ([]Movie, error)
	MovieById(id int) (Movie, error)
	CreateMovie(title, genre, releaseDate string, duration int) (movieId int, err error)
	UpdateMovie(id int, title, genre, releaseDate string, duration int) error
	DeleteMovie(id int) error
	MovieExists(id int) (bool, error)
}

type Service struct {
	R repository
}

func New(r repository) Service {
	return Service{R: r}
}

func (s Service) Movies() ([]Movie, error) {
	date := time.Now().Format("2006-01-02")
	halls, err := s.R.Movies(date)
	if err != nil {
		return []Movie{}, ErrInternalError
	}
	return halls, nil
}

func (s Service) MovieById(id int) (Movie, error) {
	hall, err := s.R.MovieById(id)
	if errors.Is(err, ErrMovieNotFound) {
		return Movie{}, ErrMovieNotFound
	}
	if err != nil {
		return Movie{}, ErrInternalError
	}
	return hall, nil
}

func (s Service) CreateMovie(title, genre, releaseDate string, duration int) (movieId int, err error) {
	id, err := s.R.CreateMovie(title, genre, releaseDate, duration)
	if err != nil {
		return 0, ErrInternalError
	}
	return id, nil
}

func (s Service) UpdateMovie(id int, title, genre, releaseDate string, duration int) error {
	ok, err := s.R.MovieExists(id)
	if err != nil {
		return ErrInternalError
	}
	if !ok {
		return ErrMovieNotFound
	}
	return s.R.UpdateMovie(id, title, genre, releaseDate, duration)
}

func (s Service) DeleteMovie(id int) error {
	ok, err := s.R.MovieExists(id)
	if err != nil {
		return ErrInternalError
	}
	if !ok {
		return ErrMovieNotFound
	}
	return s.R.DeleteMovie(id)
}
