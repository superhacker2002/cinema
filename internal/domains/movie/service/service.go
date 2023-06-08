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

const AdminRole = "admin"

type Movie struct {
	Id          int
	Title       string
	Genre       string
	ReleaseDate string
	Duration    int
}

func NewMovieEntity(id int, title, genre string, releaseDate time.Time, duration int) Movie {
	return Movie{
		Id:          id,
		Title:       title,
		Genre:       genre,
		ReleaseDate: releaseDate.Format("2006-01-02"),
		Duration:    duration,
	}
}

type repository interface {
	Movies(date string) ([]Movie, error)
	MovieById(id int) (Movie, error)
	CreateMovie(title, genre, releaseDate string, duration int) (movieId int, err error)
	UpdateMovie(id int, title, genre, releaseDate string, duration int) (bool, error)
	DeleteMovie(id int) (bool, error)
	WatchedMovies(userId int) ([]Movie, error)
	UserExists(id int) (bool, error)
}

type Service struct {
	r repository
}

func New(r repository) Service {
	return Service{r: r}
}

func (s Service) Movies() ([]Movie, error) {
	date := time.Now().Format("2006-01-02")
	movies, err := s.r.Movies(date)
	if err != nil {
		return []Movie{}, ErrInternalError
	}
	return movies, nil
}

func (s Service) MovieById(id int) (Movie, error) {
	movie, err := s.r.MovieById(id)
	if errors.Is(err, ErrMoviesNotFound) {
		return Movie{}, ErrMoviesNotFound
	}
	if err != nil {
		return Movie{}, ErrInternalError
	}
	return movie, nil
}

func (s Service) CreateMovie(title, genre, releaseDate string, duration int) (movieId int, err error) {
	id, err := s.r.CreateMovie(title, genre, releaseDate, duration)
	if err != nil {
		return 0, ErrInternalError
	}
	return id, nil
}

func (s Service) UpdateMovie(id int, title, genre, releaseDate string, duration int) error {
	found, err := s.r.UpdateMovie(id, title, genre, releaseDate, duration)
	if err != nil {
		return ErrInternalError
	}
	if !found {
		return ErrMoviesNotFound
	}

	return nil
}

func (s Service) DeleteMovie(id int) error {
	found, err := s.r.DeleteMovie(id)
	if err != nil {
		return ErrInternalError
	}
	if !found {
		return ErrMoviesNotFound
	}

	return nil
}

func (s Service) WatchedMovies(userId int) ([]Movie, error) {
	ok, err := s.r.UserExists(userId)
	if err != nil {
		return nil, ErrInternalError
	}

	if !ok {
		return nil, ErrUserNotFound
	}

	movies, err := s.r.WatchedMovies(userId)

	if err != nil {
		return nil, ErrInternalError
	}

	return movies, nil
}
