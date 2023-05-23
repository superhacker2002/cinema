package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/repository"
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
	ErrInternalError    = errors.New("internal server error")
	ErrInvalidDate      = errors.New("invalid date format")
	ErrInvalidSessionId = errors.New("invalid session id")
	ErrReadRequestFail  = errors.New("failed to read request body")
)

type HttpHandler struct {
	r repository.Repository
}

type Page struct {
	Offset int
	Limit  int
}

type session struct {
	Id        int    `json:"id"`
	MovieId   int    `json:"movieId"`
	StartTime string `json:"startTime"`
	Status    string `json:"status"`
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

	sessions, err := h.r.AllSessions(d, p.Offset, p.Limit)

	if errors.Is(err, repository.ErrCinemaSessionsNotFound) {
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
	err = json.NewEncoder(w).Encode(convert(sessions))
	if err != nil {
		log.Println(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}
}

func (h HttpHandler) getSessionsHandler(w http.ResponseWriter, r *http.Request) {
	hallId, err := pathVariable(r, "hallId")
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v: %d", ErrInvalidHallId, hallId), http.StatusBadRequest)
		return
	}

	d, err := date(r)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v: %s", ErrInvalidDate, d), http.StatusBadRequest)
		return
	}

	sessions, err := h.r.SessionsForHall(hallId, d)

	if errors.Is(err, repository.ErrCinemaSessionsNotFound) {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v for hall %d", err, hallId), http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(convert(sessions))
	if err != nil {
		log.Println(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}
}

func (h HttpHandler) createSessionHandler(w http.ResponseWriter, r *http.Request) {
	hallId, err := pathVariable(r, "hallId")
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v: %d", ErrInvalidHallId, hallId), http.StatusBadRequest)
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

	//id, err := h.r.CreateUser(creds.Username, creds.Password)
	//if errors.Is(err, userRepository.ErrUserExists) {
	//	log.Println(err)
	//	http.Error(w, err.Error(), http.StatusConflict)
	//	return
	//}
	//
	//if err != nil {
	//	log.Println(err)
	//	http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
	//	return
	//}
	//
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(map[string]int{"user_id": id})
	//w.WriteHeader(http.StatusOK)
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
	sessionId, err := pathVariable(r, "sessionId")
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v: %d", ErrInvalidSessionId, sessionId), http.StatusBadRequest)
		return
	}

	err = h.r.DeleteSession(sessionId)
	if errors.Is(repository.ErrCinemaSessionsNotFound, err) {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	if err != nil {
		log.Println(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
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

func convert(sessions []repository.CinemaSession) []session {
	var jsonSessions []session
	for _, s := range sessions {
		jsonSession := session{
			Id:        s.ID,
			MovieId:   s.MovieId,
			StartTime: s.StartTime,
			Status:    s.Status,
		}

		jsonSessions = append(jsonSessions, jsonSession)
	}
	return jsonSessions
}
