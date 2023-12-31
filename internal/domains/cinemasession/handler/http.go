package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/apiutils"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/cinemasession/entity"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/cinemasession/service"
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

const timestampLayout = "2006-01-02 15:04:05 MST"

var (
	ErrInvalidHallId    = errors.New("invalid hall id")
	ErrInvalidDate      = errors.New("invalid date format")
	ErrInvalidSessionId = errors.New("invalid session id")
	ErrReadRequestFail  = errors.New("failed to read request body")
	ErrInvalidOffset    = errors.New("invalid offset parameter")
	ErrInvalidLimit     = errors.New("invalid limit parameter")
)

type Service interface {
	AllSessions(date string, offset, limit int) ([]entity.CinemaSession, error)
	SessionsForHall(hallId int, date string) ([]entity.CinemaSession, error)
	CreateSession(movieId, hallId int, startTime string, price float32) (int, error)
	DeleteSession(id int) error
	UpdateSession(id, movieId, hallId int, startTime string, price float32) error
	AvailableSeats(sessionId int) ([]int, error)
}

type AccessChecker interface {
	Authenticate(next http.Handler) http.Handler
	CheckPerms(perms ...string) mux.MiddlewareFunc
}

type HttpHandler struct {
	s Service
}

type Page struct {
	offset int
	limit  int
}

type session struct {
	Id        int     `json:"id"`
	MovieId   int     `json:"movieId"`
	HallId    int     `json:"hallId"`
	StartTime string  `json:"startTime"`
	Price     float32 `json:"price"`
	Status    string  `json:"status"`
}

func New(s Service) HttpHandler {
	return HttpHandler{
		s: s,
	}
}

func (h HttpHandler) SetRoutes(router *mux.Router, a AccessChecker) {
	userRouter := router.PathPrefix("/cinema-sessions").Subrouter()
	userRouter.Use(a.Authenticate)

	userRouter.HandleFunc("/", h.getAllSessionsHandler).Methods("GET")
	userRouter.HandleFunc("/{hallId}", h.getSessionsHandler).Methods("GET")
	userRouter.HandleFunc("/{sessionId}/seats", h.availableSeatsHandler).Methods("GET")

	adminRouter := router.PathPrefix("/cinema-sessions").Subrouter()
	adminRouter.Use(a.Authenticate)
	adminRouter.Use(a.CheckPerms(service.AdminRole))

	adminRouter.HandleFunc("/{sessionId}", h.updateSessionHandler).Methods("PUT")
	adminRouter.HandleFunc("/{sessionId}", h.deleteSessionHandler).Methods("DELETE")
	adminRouter.HandleFunc("/", h.createSessionHandler).Methods("POST")
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

	if errors.Is(err, service.ErrCinemaSessionsNotFound) {
		http.Error(w, fmt.Sprintf("%v for all halls", err), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, entitiesToDTO(sessions), http.StatusOK)
}

func (h HttpHandler) getSessionsHandler(w http.ResponseWriter, r *http.Request) {
	hallId, err := apiutils.IntPathParam(r, "hallId")
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, entitiesToDTO(sessions), http.StatusOK)
}

func (h HttpHandler) createSessionHandler(w http.ResponseWriter, r *http.Request) {
	type sessionInfo struct {
		MovieId   int     `json:"movieId"`
		HallId    int     `json:"hallId"`
		StartTime string  `json:"startTime"`
		Price     float32 `json:"price"`
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

	id, err := h.s.CreateSession(session.MovieId, session.HallId, session.StartTime, session.Price)

	if errors.Is(err, service.ErrHallIsBusy) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if errors.Is(err, service.ErrHallNotFound) || errors.Is(err, service.ErrMovieNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	apiutils.WriteResponse(w, map[string]int{"sessionId": id}, http.StatusCreated)
}

func (h HttpHandler) updateSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionId, err := apiutils.IntPathParam(r, "sessionId")
	if err != nil {
		log.Println(err)
		http.Error(w, ErrInvalidSessionId.Error(), http.StatusBadRequest)
		return
	}

	type sessionInfo struct {
		MovieId   int     `json:"movieId"`
		HallId    int     `json:"hallId"`
		StartTime string  `json:"startTime"`
		Price     float32 `json:"price"`
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

	err = h.s.UpdateSession(sessionId, session.MovieId, session.HallId, session.StartTime, session.Price)

	if errors.Is(err, service.ErrHallIsBusy) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if errors.Is(err, service.ErrHallNotFound) || errors.Is(err, service.ErrMovieNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) deleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionId, err := apiutils.IntPathParam(r, "sessionId")
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpHandler) availableSeatsHandler(w http.ResponseWriter, r *http.Request) {
	sessionId, err := apiutils.IntPathParam(r, "sessionId")
	if err != nil {
		log.Println(err)
		http.Error(w, ErrInvalidSessionId.Error(), http.StatusBadRequest)
		return
	}

	seats, err := h.s.AvailableSeats(sessionId)
	if errors.Is(service.ErrCinemaSessionsNotFound, err) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, seats, http.StatusOK)

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
		return p, ErrInvalidOffset
	}
	if p.limit, err = strconv.Atoi(limitStr); err != nil {
		return p, ErrInvalidLimit
	}

	if p.offset < 0 {
		return p, ErrInvalidOffset
	}
	if p.limit < 0 {
		return p, ErrInvalidLimit
	}

	return p, nil
}

func date(r *http.Request) (string, error) {
	date := r.URL.Query().Get("date")
	const dateLayout = "2006-01-02"
	if date == "" {
		return time.Now().Format(dateLayout), nil
	}
	_, err := time.Parse(dateLayout, date)
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
			HallId:    s.HallId,
			StartTime: s.StartTime.Format(timestampLayout),
			Price:     s.Price,
			Status:    s.Status,
		})
	}
	return DTOSessions
}
