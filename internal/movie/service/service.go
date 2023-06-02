package service

import (
	"errors"
	"time"
)

var (
	ErrMoviesNotFound = errors.New("movies not found")
	ErrInternalError  = errors.New("internal server error")
	ErrUserNotFound   = errors.New("user not found")
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
	WatchedMovies(userId int) ([]Movie, error)
	UserExists(id int) (bool, error)
}

type Service struct {
	R repository
}

func New(r repository) Service {
	return Service{R: r}
}

func (s Service) Movies() ([]Movie, error) {
	date := time.Now().Format("2006-01-02")
	movies, err := s.R.Movies(date)
	if err != nil {
		return []Movie{}, ErrInternalError
	}
	return movies, nil
}

func (s Service) MovieById(id int) (Movie, error) {
	movie, err := s.R.MovieById(id)
	if errors.Is(err, ErrMoviesNotFound) {
		return Movie{}, ErrMoviesNotFound
	}
	if err != nil {
		return Movie{}, ErrInternalError
	}
	return movie, nil
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
		return ErrMoviesNotFound
	}
	if err = s.R.UpdateMovie(id, title, genre, releaseDate, duration); err != nil {
		return ErrInternalError
	}
	return nil
}

func (s Service) DeleteMovie(id int) error {
	ok, err := s.R.MovieExists(id)
	if err != nil {
		return ErrInternalError
	}
	if !ok {
		return ErrMoviesNotFound
	}
	if err = s.R.DeleteMovie(id); err != nil {
		return ErrInternalError
	}
	return nil
}

func (s Service) WatchedMovies(userId int) ([]Movie, error) {
	ok, err := s.R.UserExists(userId)
	if err != nil {
		return nil, ErrInternalError
	}
	if !ok {
		return nil, ErrUserNotFound
	}
	movies, err := s.R.WatchedMovies(userId)
	if err != nil {
		return nil, ErrInternalError
	}
	return movies, nil
}
