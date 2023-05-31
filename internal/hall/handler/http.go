package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/hall/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/apiutils"
	"github.com/gorilla/mux"
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

type Repository interface {
	Halls() ([]service.Hall, error)
	HallById(id int) (service.Hall, error)
	CreateHall(name string, capacity int) (hallId int, err error)
	UpdateHall(id int, name string, capacity int) error
	DeleteHall(id int) error
	//UpdateHallAvailability(id int, available bool) error
}

type HTTPHandler struct {
	repository Repository
}

func New(router *mux.Router, repository Repository) HTTPHandler {
	handler := HTTPHandler{repository: repository}
	handler.setRoutes(router)

	return handler
}

func (h HTTPHandler) setRoutes(router *mux.Router) {
	s := router.PathPrefix("/halls").Subrouter()
	s.HandleFunc("/", h.getHallsHandler).Methods(http.MethodGet)
	s.HandleFunc("/", h.createHallHandler).Methods(http.MethodPost)
	s.HandleFunc("/{hallId}/", h.getHallHandler).Methods(http.MethodGet)
	s.HandleFunc("/{hallId}/", h.updateHallHandler).Methods(http.MethodPut)
	s.HandleFunc("/{hallId}/", h.deleteHallHandler).Methods(http.MethodDelete)
	s.HandleFunc("/update-availability", h.updateAvailabilityHandler).Methods(http.MethodPut)
}

func (h HTTPHandler) getHallsHandler(w http.ResponseWriter, _ *http.Request) {
	cinemaHalls, err := h.repository.Halls()
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

	id, err := h.repository.CreateHall(hall.Name, hall.Capacity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, map[string]int{"user_id": id}, http.StatusCreated)
}

func (h HTTPHandler) getHallHandler(w http.ResponseWriter, r *http.Request) {
	hallID, err := apiutils.IntPathParam(r, "hallId")
	if err != nil {
		http.Error(w, ErrInvalidHallId.Error(), http.StatusBadRequest)
		return
	}

	hall, err := h.repository.HallById(hallID)
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

	err = h.repository.UpdateHall(hallID, hall.Name, hall.Capacity)
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

func (h HTTPHandler) deleteHallHandler(w http.ResponseWriter, r *http.Request) {
	hallID, err := apiutils.IntPathParam(r, "hallID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.repository.DeleteCinemaHall(hallID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Cinema hall with ID %d deleted", hallID)
	apiutils.WriteResponse(w, message, http.StatusOK)
}

func (h HTTPHandler) updateAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	var update struct {
		Data struct {
			HallID    int  `json:"hallId"`
			Available bool `json:"available"`
		} `json:"data"`
	}

	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if update.Data.HallID <= 0 {
		http.Error(w, "Invalid HallID", http.StatusBadRequest)
		return
	}

	err = h.repository.UpdateHallAvailability(update.Data.HallID, update.Data.Available)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Updated availability for cinema hall with ID %d", update.Data.HallID)
	apiutils.WriteResponse(w, message, http.StatusOK)
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
