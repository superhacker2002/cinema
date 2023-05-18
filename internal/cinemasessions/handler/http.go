package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/repository"
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
	r repository.Repository
}

func New(router *mux.Router, repository repository.Repository) HttpHandler {
	handler := HttpHandler{r: repository}
	handler.setRoutes(router)

	return handler
}

func (h HttpHandler) setRoutes(router *mux.Router) {
	s := router.PathPrefix("/cinema-sessions").Subrouter()
	s.HandleFunc("/", h.getAllSessionsHandler).Methods("GET")
	s.HandleFunc("/{hallId}", h.getSessionsHandler).Methods("GET")
	//s.HandleFunc("/", h.createSessionHandler).Methods("POST")
	//s.HandleFunc("/{sessionId}/", h.getSessionHandler).Methods("GET")
	//s.HandleFunc("/{sessionId}/", h.updateSessionHandler).Methods("PUT")
	//s.HandleFunc("/{sessionId}/", h.deleteSessionHandler).Methods("DELETE")
}

func (h HttpHandler) getAllSessionsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: add getting offset and limit from URL
	date := r.URL.Query().Get("date")
	sessions, err := h.r.AllSessions(date, 0, 10)

	if errors.Is(err, repository.ErrCinemaSessionsNotFound) {
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
	date := r.URL.Query().Get("date")
	sessions, err := h.r.SessionsForHall(hallId, date, 0, 10)

	if errors.Is(err, repository.ErrCinemaSessionsNotFound) {
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

func addTime(r *http.Request) string {
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
