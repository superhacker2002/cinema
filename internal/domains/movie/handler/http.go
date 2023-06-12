package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/apiutils"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/movie/service"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
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
	ReleaseDate string `json:"releaseDate"`
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

type AccessChecker interface {
	Authenticate(next http.Handler) http.Handler
	CheckPerms(perms ...string) mux.MiddlewareFunc
}

type HttpHandler struct {
	s Service
}

func New(s Service) HttpHandler {
	return HttpHandler{
		s: s,
	}
}

func (h HttpHandler) SetRoutes(router *mux.Router, a AccessChecker) {
	userRouter := router.PathPrefix("/movies").Subrouter()
	userRouter.Use(a.Authenticate)

	userRouter.HandleFunc("/", h.getMoviesHandler).Methods(http.MethodGet)
	userRouter.HandleFunc("/{movieId}", h.getMovieHandler).Methods(http.MethodGet)
	userRouter.HandleFunc("/watched/{userId}", h.watchedMoviesHandler).Methods(http.MethodGet)

	adminRouter := router.PathPrefix("/movies").Subrouter()
	adminRouter.Use(a.Authenticate)
	adminRouter.Use(a.CheckPerms(service.AdminRole))

	adminRouter.HandleFunc("/", h.createMovieHandler).Methods(http.MethodPost)
	adminRouter.HandleFunc("/{movieId}", h.updateMovieHandler).Methods(http.MethodPut)
	adminRouter.HandleFunc("/{movieId}", h.deleteMovieHandler).Methods(http.MethodDelete)
}

func (h HttpHandler) getMoviesHandler(w http.ResponseWriter, _ *http.Request) {
	cinemaHalls, err := h.s.Movies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, entitiesToDTO(cinemaHalls), http.StatusOK)
}

func (h HttpHandler) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	type movieInfo struct {
		Title       string `json:"title"`
		Genre       string `json:"genre"`
		ReleaseDate string `json:"releaseDate"`
		Duration    int    `json:"duration"`
	}

	var movie movieInfo

	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.s.CreateMovie(movie.Title, movie.Genre, movie.ReleaseDate, movie.Duration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, map[string]int{"movieId": id}, http.StatusCreated)
}

func (h HttpHandler) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieId, err := apiutils.IntPathParam(r, "movieId")
	if err != nil {
		http.Error(w, ErrInvalidMovieId.Error(), http.StatusBadRequest)
		return
	}

	movie, err := h.s.MovieById(movieId)
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

func (h HttpHandler) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieId, err := apiutils.IntPathParam(r, "movieId")
	if err != nil {
		http.Error(w, ErrInvalidMovieId.Error(), http.StatusBadRequest)
		return
	}

	type movieInfo struct {
		Title       string `json:"title"`
		Genre       string `json:"genre"`
		ReleaseDate string `json:"releaseDate"`
		Duration    int    `json:"duration"`
	}

	var movie movieInfo

	err = json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.UpdateMovie(movieId, movie.Title, movie.Genre, movie.ReleaseDate, movie.Duration)
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

func (h HttpHandler) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieId, err := apiutils.IntPathParam(r, "movieId")
	if err != nil {
		http.Error(w, ErrInvalidMovieId.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.DeleteMovie(movieId)
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

func (h HttpHandler) watchedMoviesHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := apiutils.IntPathParam(r, "userId")
	if err != nil {
		http.Error(w, ErrInvalidUserId.Error(), http.StatusBadRequest)
		return
	}
	movies, err := h.s.WatchedMovies(userId)
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
