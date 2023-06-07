package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/apiutils"
	ticketServ "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/ticket/service"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	ErrReadRequestFail = errors.New("failed to read request body")
)

type service interface {
	BuyTicket(sessionId, userId, seatNum int) (string, error)
}

type accessChecker interface {
	Authenticate(next http.Handler) http.Handler
	CheckPerms(next http.Handler, perms ...string) http.Handler
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
	SeatNumber int `json:"seatNumber"`
}

func (h HttpHandler) SetRoutes(router *mux.Router, a accessChecker) {
	s := router.PathPrefix("/tickets").Subrouter()
	s.Use(a.Authenticate)
	s.HandleFunc("/", h.createTicket).Methods(http.MethodPost)
}

func (h HttpHandler) createTicket(w http.ResponseWriter, r *http.Request) {
	var t ticket
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)

	ticketPath, err := h.s.BuyTicket(t.SessionId, userID, t.SeatNumber)
	if errors.Is(err, ticketServ.ErrCinemaSessionsNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if errors.Is(err, ticketServ.ErrTicketExists) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, map[string]string{"ticketPath": ticketPath}, http.StatusCreated)
}
