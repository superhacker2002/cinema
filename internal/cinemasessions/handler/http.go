package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/entity"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type repository interface {
	SessionsForHall(hallId int, date string) ([]entity.CinemaSession, error)
	AllSessions(date string, offset, limit int) ([]entity.CinemaSession, error)
}

var (
	ErrInvalidHallId    = errors.New("invalid hall id")
	ErrInternalError    = errors.New("internal server error")
	ErrInvalidDate      = errors.New("invalid date format")
	ErrInvalidSessionId = errors.New("invalid session id")
)

type HttpHandler struct {
	r repository
}

type Page struct {
	Offset int
	Limit  int
}

type session struct {
	Id        int       `json:"id"`
	MovieId   int       `json:"movieId"`
	HallId    string    `json:"hallId,omitempty"`
	StartTime time.Time `json:"startTime"`
	Status    string    `json:"status"`
}

func New(router *mux.Router, repository repository) HttpHandler {
	handler := HttpHandler{r: repository}
	handler.setRoutes(router)

	return handler
}

func (h HttpHandler) setRoutes(router *mux.Router) {
	s := router.PathPrefix("/cinema-sessions").Subrouter()
	s.HandleFunc("/", h.getAllSessionsHandler).Methods("GET")
	s.HandleFunc("/{hallId}", h.getSessionsHandler).Methods("GET")
}

func (h HttpHandler) getAllSessionsHandler(w http.ResponseWriter, r *http.Request) {
	d, err := date(r)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v: %s", ErrInvalidDate, d), http.StatusBadRequest)
		return
	}

	p, err := page(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessions, err := h.r.AllSessions(d, p.Offset, p.Limit)

	if errors.Is(err, entity.ErrCinemaSessionsNotFound) {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v for all halls", err), http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(entitiesToDTO(sessions))
	if err != nil {
		log.Println(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}
}

func (h HttpHandler) getSessionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hallIdStr := vars["hallId"]
	hallId, err := strconv.Atoi(hallIdStr)
	if err != nil || hallId <= 0 {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v: %s", ErrInvalidHallId, hallIdStr), http.StatusBadRequest)
		return
	}

	d, err := date(r)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v: %s", ErrInvalidDate, d), http.StatusBadRequest)
		return
	}

	sessions, err := h.r.SessionsForHall(hallId, d)

	if errors.Is(err, entity.ErrCinemaSessionsNotFound) {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v for hall %s", err, hallIdStr), http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(entitiesToDTO(sessions))
	if err != nil {
		log.Println(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
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

func page(r *http.Request) (Page, error) {
	const (
		defaultOffset = 0
		defaultLimit  = 10
	)
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")
	var p Page
	if offsetStr == "" || limitStr == "" {
		p.Offset = defaultOffset
		p.Limit = defaultLimit
		log.Println("missing offset or limit, default values are used")
		return p, nil
	}

	var err error
	if p.Offset, err = strconv.Atoi(offsetStr); err != nil {
		return p, err
	}
	if p.Limit, err = strconv.Atoi(limitStr); err != nil {
		return p, err
	}

	if p.Offset < 0 {
		return p, errors.New(fmt.Sprintf("invalid offset parameter: %d", p.Offset))
	}
	if p.Limit < 0 {
		return p, errors.New(fmt.Sprintf("invalid limit parameter: %d", p.Offset))
	}

	return p, nil
}

func date(r *http.Request) (string, error) {
	date := r.URL.Query().Get("date")
	layout := "2006-01-02"
	if date == "" {
		return time.Now().Format(layout), nil
	}
	_, err := time.Parse(layout, date)
	if err != nil {
		return "", err
	}
	return date, nil
}

func entitiesToDTO(sessions []entity.CinemaSession) []session {
	DTOSessions := make([]session, len(sessions))
	for _, s := range sessions {
		DTOSessions = append(DTOSessions, session{
			Id:        s.Id,
			MovieId:   s.MovieId,
			StartTime: s.StartTime,
			Status:    s.Status,
		})
	}
	return DTOSessions
}
