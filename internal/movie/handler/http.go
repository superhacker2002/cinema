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
	ErrInvalidMovie    = errors.New("invalid movie id")
)

type cinemaHall struct {
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
	UpdateMovie(id int, name string, capacity int) error
	DeleteMovie(id int) error
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
	s := router.PathPrefix("/halls").Subrouter()
	s.HandleFunc("/", h.getMoviesHandler).Methods(http.MethodGet)
	s.HandleFunc("/", h.createMovieHandler).Methods(http.MethodPost)
	s.HandleFunc("/{hallId}", h.getMovieHandler).Methods(http.MethodGet)
	s.HandleFunc("/{hallId}", h.updateMovieHandler).Methods(http.MethodPut)
	s.HandleFunc("/{hallId}", h.deleteMovieHandler).Methods(http.MethodDelete)
}

func (h HTTPHandler) getHallsHandler(w http.ResponseWriter, _ *http.Request) {
	cinemaHalls, err := h.S.Halls()
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

func (h HTTPHandler) getHallHandler(w http.ResponseWriter, r *http.Request) {
	hallID, err := apiutils.IntPathParam(r, "hallId")
	if err != nil {
		http.Error(w, ErrInvalidHallId.Error(), http.StatusBadRequest)
		return
	}

	hall, err := h.S.HallById(hallID)
	if errors.Is(err, service.ErrHallNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, entityToDTO(hall), http.StatusOK)
}

func (h HTTPHandler) updateHallHandler(w http.ResponseWriter, r *http.Request) {
	hallID, err := apiutils.IntPathParam(r, "hallId")
	if err != nil {
		http.Error(w, ErrInvalidHallId.Error(), http.StatusBadRequest)
		return
	}

	type hallInfo struct {
		Name     string `json:"name"`
		Capacity int    `json:"capacity"`
	}

	var hall hallInfo

	err = json.NewDecoder(r.Body).Decode(&hall)
	if err != nil {
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	err = h.S.UpdateHall(hallID, hall.Name, hall.Capacity)
	if errors.Is(err, service.ErrHallNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteMsg(w, "cinema hall was updated successfully\n", http.StatusOK)
}

func (h HTTPHandler) deleteHallHandler(w http.ResponseWriter, r *http.Request) {
	hallID, err := apiutils.IntPathParam(r, "hallId")
	if err != nil {
		http.Error(w, ErrInvalidHallId.Error(), http.StatusBadRequest)
		return
	}

	err = h.S.DeleteHall(hallID)
	if errors.Is(err, service.ErrHallNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteMsg(w, "cinema hall was deleted successfully\n", http.StatusOK)
}

func entitiesToDTO(halls []service.Hall) []cinemaHall {
	var DTOHalls []cinemaHall
	for _, hall := range halls {
		DTOHalls = append(DTOHalls, entityToDTO(hall))
	}
	return DTOHalls
}

func entityToDTO(hall service.Hall) cinemaHall {
	return cinemaHall{
		ID:       hall.Id,
		Name:     hall.Name,
		Capacity: hall.Capacity,
	}
}
