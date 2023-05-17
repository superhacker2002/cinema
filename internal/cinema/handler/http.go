package handler

import (
	cinemaRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinema/repository"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrInvalidHallId = errors.New("invalid hall id provided")
	ErrInternalError = errors.New("internal server error")
)

type HttpHandler struct {
	r cinemaRepository.Repository
}

func New(router *mux.Router, repository cinemaRepository.Repository) HttpHandler {
	handler := HttpHandler{r: repository}
	handler.setRoutes(router)

	return handler
}

func (h HttpHandler) setRoutes(router *mux.Router) {
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
	s.HandleFunc("/", h.getAllSessionsHandler).Methods("GET")
	s.HandleFunc("/{hallId}", h.getSessionsHandler).Methods("GET")
	//s.HandleFunc("/", h.createSessionHandler).Methods("POST")
	//s.HandleFunc("/{sessionId}/", h.getSessionHandler).Methods("GET")
	//s.HandleFunc("/{sessionId}/", h.updateSessionHandler).Methods("PUT")
	//s.HandleFunc("/{sessionId}/", h.deleteSessionHandler).Methods("DELETE")

	s = router.PathPrefix("/tickets").Subrouter()
	s.HandleFunc("/", h.createTicketHandler).Methods("POST")
	s.HandleFunc("/{ticketId}/", h.getTicketHandler).Methods("GET")
	s.HandleFunc("/{userId}/", h.getUserTicketsHandler).Methods("GET")
}

func (h HttpHandler) getHallsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return all halls
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) createHallHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: create new hall (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) getHallHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return hall by id
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) updateHallHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: update hall by id (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) deleteHallHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: delete hall by id (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) getMoviesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return all movies
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: create new movie (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return movie by id
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: update movie by id (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: delete movie by id (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) watchedMoviesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return list of watched movies by user
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) getAllSessionsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: add getting offset and limit from URL
	sessions, err := h.r.AllSessions(date(r), 0, 10)

	if errors.Is(err, cinemaRepository.ErrCinemaSessionsNotFound) {
		log.Println(err)
		http.Error(w, err.Error()+"for all halls", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

func (h HttpHandler) getSessionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hallIdStr := vars["hallId"]
	hallId, err := strconv.Atoi(hallIdStr)
	if err != nil || hallId <= 0 {
		log.Println(err)
		http.Error(w, ErrInvalidHallId.Error()+": "+hallIdStr, http.StatusBadRequest)
		return
	}

	// TODO: add getting offset and limit from URL
	sessions, err := h.r.SessionsForHall(hallId, date(r), 0, 10)

	if errors.Is(err, cinemaRepository.ErrCinemaSessionsNotFound) {
		log.Println(err)
		http.Error(w, err.Error()+" for hall "+hallIdStr, http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

func date(r *http.Request) string {
	date := r.URL.Query().Get("date")
	if date == "" {
		return time.Now().Format("2006-01-02 15:04:05")
	} else {
		return date + " 00:00:00"
	}
}

func (h HttpHandler) createSessionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: create new cinema session (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) getSessionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return cinema session by id
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) updateSessionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: update cinema session by id (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) deleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: delete cinema session by id (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) createTicketHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: create new ticket
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) getTicketHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return ticket by id
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) getUserTicketsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return list of tickets purchased by user
	w.WriteHeader(http.StatusOK)
}
