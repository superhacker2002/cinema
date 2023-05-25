package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/entity"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/service"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrInvalidHallId    = errors.New("invalid hall id")
	ErrInvalidDate      = errors.New("invalid date format")
	ErrInvalidSessionId = errors.New("invalid session id")
	ErrReadRequestFail  = errors.New("failed to read request body")
)

type Service interface {
	AllSessions(date string, offset, limit int) ([]entity.CinemaSession, error)
	SessionsForHall(hallId int, date string) ([]entity.CinemaSession, error)
	CreateSession(movieId, hallId int, startTime string, price float32) (int, error)
	DeleteSession(id int) error
}

type HttpHandler struct {
	s Service
}

type Page struct {
	offset int
	limit  int
}

type session struct {
	Id        int       `json:"id"`
	MovieId   int       `json:"movieId"`
	HallId    string    `json:"hallId,omitempty"`
	StartTime time.Time `json:"startTime"`
	Status    string    `json:"status"`
}

func New(router *mux.Router, s Service) HttpHandler {
	handler := HttpHandler{s: s}
	handler.setRoutes(router)

	return handler
}

func (h HttpHandler) setRoutes(router *mux.Router) {
	s := router.PathPrefix("/cinema-sessions").Subrouter()
	s.HandleFunc("/", h.getAllSessionsHandler).Methods("GET")
	s.HandleFunc("/{sessionId}", h.updateSessionHandler).Methods("PUT")
	s.HandleFunc("/{sessionId}", h.deleteSessionHandler).Methods("DELETE")
	s.HandleFunc("/{hallId}", h.getSessionsHandler).Methods("GET")
	s.HandleFunc("/{hallId}", h.createSessionHandler).Methods("POST")
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

	sessions, err := h.s.AllSessions(d, p.offset, p.limit)
	log.Println(len(sessions))

	if errors.Is(err, service.ErrCinemaSessionsNotFound) {
		http.Error(w, fmt.Sprintf("%v for all halls", err), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, service.ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(entitiesToDTO(sessions))
	if err != nil {
		log.Println(err)
		http.Error(w, service.ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}
}

func (h HttpHandler) getSessionsHandler(w http.ResponseWriter, r *http.Request) {
	hallId, err := pathVariable(r, "hallId")
	if err != nil {
		log.Println(err)
		http.Error(w, ErrInvalidHallId.Error(), http.StatusBadRequest)
		return
	}

	d, err := date(r)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v: %s", ErrInvalidDate, d), http.StatusBadRequest)
		return
	}

	sessions, err := h.s.SessionsForHall(hallId, d)

	if errors.Is(err, service.ErrCinemaSessionsNotFound) {
		http.Error(w, fmt.Sprintf("%v for hall %d", err, hallId), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, service.ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(entitiesToDTO(sessions))
	if err != nil {
		log.Println(err)
		http.Error(w, service.ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}
}

func (h HttpHandler) createSessionHandler(w http.ResponseWriter, r *http.Request) {
	hallId, err := pathVariable(r, "hallId")
	if err != nil {
		log.Println(err)
		http.Error(w, ErrInvalidHallId.Error(), http.StatusBadRequest)
		return
	}
	type sessionInfo struct {
		MovieId  int     `json:"movieId"`
		StarTime string  `json:"starTime"`
		Price    float32 `json:"price"`
	}
	var session sessionInfo

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body, &session); err != nil {
		log.Println(err)
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.s.CreateSession(session.MovieId, hallId, session.StarTime, session.Price)

	if errors.Is(err, service.ErrHallIsBusy) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if err != nil {
		http.Error(w, service.ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"session_id": id})
	if err != nil {
		log.Println(err)
		http.Error(w, service.ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}
}

func (h HttpHandler) updateSessionHandler(w http.ResponseWriter, r *http.Request) {
	//sessionId, err := pathVariable(r, "sessionId")
	//if err != nil {
	//	log.Println(err)
	//	http.Error(w, ErrInvalidSessionId.Error(), http.StatusBadRequest)
	//	return
	//}

	type sessionInfo struct {
		MovieId   int       `json:"movieId"`
		HallId    int       `json:"hallId"`
		StartTime time.Time `json:"startTime"`
		Price     float32   `json:"price"`
	}

	var session sessionInfo

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body, &session); err != nil {
		log.Println(err)
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) deleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionId, err := pathVariable(r, "sessionId")
	if err != nil {
		log.Println(err)
		http.Error(w, ErrInvalidSessionId.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.DeleteSession(sessionId)
	if errors.Is(service.ErrCinemaSessionsNotFound, err) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, service.ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte("session was deleted successfully"))
	if err != nil {
		log.Println(err)
		http.Error(w, service.ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}
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
		p.offset = defaultOffset
		p.limit = defaultLimit
		log.Println("missing offset or limit, default values are used")
		return p, nil
	}

	var err error
	if p.offset, err = strconv.Atoi(offsetStr); err != nil {
		return p, err
	}
	if p.limit, err = strconv.Atoi(limitStr); err != nil {
		return p, err
	}

	if p.offset < 0 {
		return p, errors.New(fmt.Sprintf("invalid offset parameter: %d", p.offset))
	}
	if p.limit < 0 {
		return p, errors.New(fmt.Sprintf("invalid limit parameter: %d", p.offset))
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
		return date, err
	}
	return date, nil
}

func entitiesToDTO(sessions []entity.CinemaSession) []session {
	var DTOSessions []session
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

func pathVariable(r *http.Request, varName string) (int, error) {
	vars := mux.Vars(r)
	varStr := vars[varName]
	varInt, err := strconv.Atoi(varStr)
	if err != nil {
		return 0, err
	}
	if varInt <= 0 {
		return 0, errors.New("parameter is less than zero")
	}
	return varInt, nil
}
