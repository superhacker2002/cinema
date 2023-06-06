package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/apiutils"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	ErrReadRequestFail = errors.New("failed to read request body")
)

type service interface {
	BuyTicket(sessionId, userId, seatNum int) (string, error)
}

type accessChecker interface {
	Authorize(next http.Handler) http.Handler
	CheckPerms(next http.Handler, perms []string) http.Handler
}

func New(s service) HttpHandler {
	return HttpHandler{
		s: s,
	}
}

type HttpHandler struct {
	s service
}

type ticket struct {
	SessionId  int `json:"sessionId"`
	UserId     int `json:"userId"`
	SeatNumber int `json:"seatNumber"`
}

func (h HttpHandler) SetRoutes(router *mux.Router, a accessChecker) {
	s := router.PathPrefix("/tickets").Subrouter()
	s.Use(a.Authorize)
	s.HandleFunc("/", h.createTicket).Methods(http.MethodPost)
}

func (h HttpHandler) createTicket(w http.ResponseWriter, r *http.Request) {
	var t ticket
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	ticketPath, err := h.s.BuyTicket(t.SessionId, t.UserId, t.SeatNumber)
	if err != nil {
		log.Fatal(err)
	}

	apiutils.WriteResponse(w, map[string]string{"ticketPath": ticketPath}, http.StatusCreated)
}
