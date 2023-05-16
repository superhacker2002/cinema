package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var ErrInvalidHallId = errors.New("invalid hall id provided")

type repository interface {
	SessionsForHall(hallId int) error
}

type httpHandler struct {
	r repository
}

func New(router *mux.Router, repository repository) httpHandler {
	handler := httpHandler{r: repository}
	handler.setRoutes(router)

	return handler
}

func (h httpHandler) setRoutes(router *mux.Router) {
	s := router.PathPrefix("/halls").Subrouter()
	s.HandleFunc("/", h.getHallsHandler).Methods("GET")
	s.HandleFunc("/", h.createHallHandler).Methods("POST")
	s.HandleFunc("/{hallId}/", h.getHallHandler).Methods("GET")
	s.HandleFunc("/{hallId}/", h.updateHallHandler).Methods("PUT")
	s.HandleFunc("/{hallId}/", h.deleteHallHandler).Methods("DELETE")

	s = router.PathPrefix("/movies").Subrouter()
	s.HandleFunc("/", h.getMoviesHandler).Methods("GET")
	s.HandleFunc("/", h.createMovieHandler).Methods("POST")
	s.HandleFunc("/{movieId}/", h.getMovieHandler).Methods("GET")
	s.HandleFunc("/{movieId}/", h.updateMovieHandler).Methods("PUT")
	s.HandleFunc("/{movieId}/", h.deleteMovieHandler).Methods("DELETE")
	s.HandleFunc("/watched/{userId}/", h.watchedMoviesHandler).Methods("GET")

	s = router.PathPrefix("/cinema-sessions").Subrouter()
	s.HandleFunc("/{hallId}/", h.getSessionsHandler).Methods("GET")
	//s.HandleFunc("/", h.getSessionsHandler).Methods("GET")
	s.HandleFunc("/", h.createSessionHandler).Methods("POST")
	s.HandleFunc("/{sessionId}/", h.getSessionHandler).Methods("GET")
	s.HandleFunc("/{sessionId}/", h.updateSessionHandler).Methods("PUT")
	s.HandleFunc("/{sessionId}/", h.deleteSessionHandler).Methods("DELETE")

	s = router.PathPrefix("/tickets").Subrouter()
	s.HandleFunc("/", h.createTicketHandler).Methods("POST")
	s.HandleFunc("/{ticketId}/", h.getTicketHandler).Methods("GET")
	s.HandleFunc("/{userId}/", h.getUserTicketsHandler).Methods("GET")
}

func (h httpHandler) getHallsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return all halls
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) createHallHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: create new hall (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) getHallHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return hall by id
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) updateHallHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: update hall by id (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) deleteHallHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: delete hall by id (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) getMoviesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return all movies
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: create new movie (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return movie by id
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: update movie by id (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: delete movie by id (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) watchedMoviesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return list of watched movies by user
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) getSessionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hallIdStr := vars["hallId"]
	hallId, err := strconv.Atoi(hallIdStr)
	if err != nil {
		log.Println(err)
		http.Error(w, ErrInvalidHallId.Error(), http.StatusBadRequest)
		return
	}

	if err = h.r.SessionsForHall(hallId); err != nil {
		log.Println(err)
		http.Error(w, "failed to get cinema sessions: "+ErrInvalidHallId.Error(), http.StatusBadRequest)
		return
	}
	// TODO: return today's sessions for the particular hall
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) createSessionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: create new cinema session (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) getSessionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return cinema session by id
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) updateSessionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: update cinema session by id (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) deleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: delete cinema session by id (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) createTicketHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: create new ticket
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) getTicketHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return ticket by id
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) getUserTicketsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return list of tickets purchased by user
	w.WriteHeader(http.StatusOK)
}
