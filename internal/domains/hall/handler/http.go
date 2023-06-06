package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/apiutils"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/hall/service"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	ErrReadRequestFail = errors.New("failed to read request")
	ErrInvalidHallId   = errors.New("invalid hall id")
)

type cinemaHall struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

type Service interface {
	Halls() ([]service.Hall, error)
	HallById(id int) (service.Hall, error)
	CreateHall(name string, capacity int) (hallId int, err error)
	UpdateHall(id int, name string, capacity int) (err error)
	DeleteHall(id int) error
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
	s.HandleFunc("/", h.getHallsHandler).Methods(http.MethodGet)
	s.HandleFunc("/", h.createHallHandler).Methods(http.MethodPost)
	s.HandleFunc("/{hallId}", h.getHallHandler).Methods(http.MethodGet)
	s.HandleFunc("/{hallId}", h.updateHallHandler).Methods(http.MethodPut)
	s.HandleFunc("/{hallId}", h.deleteHallHandler).Methods(http.MethodDelete)
}

func (h HTTPHandler) getHallsHandler(w http.ResponseWriter, _ *http.Request) {
	cinemaHalls, err := h.S.Halls()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, entitiesToDTO(cinemaHalls), http.StatusOK)
}

func (h HTTPHandler) createHallHandler(w http.ResponseWriter, r *http.Request) {
	type hallInfo struct {
		Name     string `json:"name"`
		Capacity int    `json:"capacity"`
	}

	var hall hallInfo

	err := json.NewDecoder(r.Body).Decode(&hall)
	if err != nil {
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.S.CreateHall(hall.Name, hall.Capacity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, map[string]int{"hallId": id}, http.StatusCreated)
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

	w.WriteHeader(http.StatusOK)
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

	w.WriteHeader(http.StatusNoContent)
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
