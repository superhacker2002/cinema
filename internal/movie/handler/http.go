package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/movie/service"
	"encoding/json"
	"errors"
	"net/http"

	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/apiutils"
	"github.com/gorilla/mux"
)

var (
	ErrReadRequestFail = errors.New("failed to read request")
	ErrInvalidMovieId  = errors.New("invalid movie id")
	ErrInvalidUserId   = errors.New("invalid user id")
)

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Genre       string `json:"genre"`
	ReleaseDate string `json:"release_date"`
	Duration    int    `json:"duration"`
}

type Service interface {
	Movies() ([]service.Movie, error)
	MovieById(id int) (service.Movie, error)
	CreateMovie(title, genre, releaseDate string, duration int) (movieId int, err error)
	UpdateMovie(id int, title, genre, releaseDate string, duration int) error
	DeleteMovie(id int) error
	WatchedMovies(userId int) ([]service.Movie, error)
}

type HTTPHandler struct {
	S Service
}

func New(router *mux.Router, s Service) HTTPHandler {
	handler := HTTPHandler{S: s}
	handler.setRoutes(router)

	return handler
}

func (h HTTPHandler) setRoutes(router *mux.Router) {
	s := router.PathPrefix("/movies").Subrouter()
	s.HandleFunc("/", h.getMoviesHandler).Methods(http.MethodGet)
	s.HandleFunc("/", h.createMovieHandler).Methods(http.MethodPost)
	s.HandleFunc("/{movieId}", h.getMovieHandler).Methods(http.MethodGet)
	s.HandleFunc("/{movieId}", h.updateMovieHandler).Methods(http.MethodPut)
	s.HandleFunc("/{movieId}", h.deleteMovieHandler).Methods(http.MethodDelete)
	s.HandleFunc("/watched/{userId}", h.watchedMoviesHandler).Methods(http.MethodGet)
}

func (h HTTPHandler) getMoviesHandler(w http.ResponseWriter, _ *http.Request) {
	cinemaHalls, err := h.S.Movies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, entitiesToDTO(cinemaHalls), http.StatusOK)
}

func (h HTTPHandler) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	type movieInfo struct {
		Title       string `json:"title"`
		Genre       string `json:"genre"`
		ReleaseDate string `json:"release_date"`
		Duration    int    `json:"duration"`
	}

	var movie movieInfo

	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.S.CreateMovie(movie.Title, movie.Genre, movie.ReleaseDate, movie.Duration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, map[string]int{"movieId": id}, http.StatusCreated)
}

func (h HTTPHandler) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieId, err := apiutils.IntPathParam(r, "movieId")
	if err != nil {
		http.Error(w, ErrInvalidMovieId.Error(), http.StatusBadRequest)
		return
	}

	movie, err := h.S.MovieById(movieId)
	if errors.Is(err, service.ErrMoviesNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, entityToDTO(movie), http.StatusOK)
}

func (h HTTPHandler) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieId, err := apiutils.IntPathParam(r, "movieId")
	if err != nil {
		http.Error(w, ErrInvalidMovieId.Error(), http.StatusBadRequest)
		return
	}

	type movieInfo struct {
		Title       string `json:"title"`
		Genre       string `json:"genre"`
		ReleaseDate string `json:"release_date"`
		Duration    int    `json:"duration"`
	}

	var movie movieInfo

	err = json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	err = h.S.UpdateMovie(movieId, movie.Title, movie.Genre, movie.ReleaseDate, movie.Duration)
	if errors.Is(err, service.ErrMoviesNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h HTTPHandler) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieId, err := apiutils.IntPathParam(r, "movieId")
	if err != nil {
		http.Error(w, ErrInvalidMovieId.Error(), http.StatusBadRequest)
		return
	}

	err = h.S.DeleteMovie(movieId)
	if errors.Is(err, service.ErrMoviesNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HTTPHandler) watchedMoviesHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := apiutils.IntPathParam(r, "userId")
	if err != nil {
		http.Error(w, ErrInvalidUserId.Error(), http.StatusBadRequest)
		return
	}
	movies, err := h.S.WatchedMovies(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, entitiesToDTO(movies), http.StatusOK)
}

func entitiesToDTO(halls []service.Movie) []Movie {
	var DTOHalls []Movie
	for _, hall := range halls {
		DTOHalls = append(DTOHalls, entityToDTO(hall))
	}
	return DTOHalls
}

func entityToDTO(hall service.Movie) Movie {
	return Movie{
		ID:          hall.Id,
		Title:       hall.Title,
		Genre:       hall.Genre,
		ReleaseDate: hall.ReleaseDate,
		Duration:    hall.Duration,
	}
}
