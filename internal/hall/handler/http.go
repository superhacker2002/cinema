package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/apiutils"
	"github.com/gorilla/mux"
)

type CinemaHall struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Capacity  int    `json:"capacity"`
	Available bool   `json:"available"`
}

type Repository interface {
	GetCinemaHalls() ([]CinemaHall, error)
	GetCinemaHallByID(id int) (*CinemaHall, error)
	CreateCinemaHall(hall *CinemaHall) error
	UpdateCinemaHall(hall *CinemaHall) error
	DeleteCinemaHall(id int) error
	UpdateHallAvailability(id int, available bool) error
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
	s.HandleFunc("/{hallID}/", h.getHallHandler).Methods(http.MethodGet)
	s.HandleFunc("/{hallID}/", h.updateHallHandler).Methods(http.MethodPut)
	s.HandleFunc("/{hallID}/", h.deleteHallHandler).Methods(http.MethodDelete)
	s.HandleFunc("/update-availability", h.updateAvailabilityHandler).Methods(http.MethodPut)
}

func (h HTTPHandler) getHallsHandler(w http.ResponseWriter, _ *http.Request) {
	cinemaHalls, err := h.repository.GetCinemaHalls()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, cinemaHalls, http.StatusOK)
}

func (h HTTPHandler) createHallHandler(w http.ResponseWriter, r *http.Request) {
	var newCinemaHall CinemaHall
	err := json.NewDecoder(r.Body).Decode(&newCinemaHall)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.repository.CreateCinemaHall(&newCinemaHall)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, newCinemaHall, http.StatusOK)
}

func (h HTTPHandler) getHallHandler(w http.ResponseWriter, r *http.Request) {
	hallID, err := apiutils.GetIntParam(r, "hallID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hall, err := h.repository.GetCinemaHallByID(hallID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	apiutils.WriteResponse(w, hall, http.StatusOK)
}

func (h HTTPHandler) updateHallHandler(w http.ResponseWriter, r *http.Request) {
	var updatedCinemaHall CinemaHall
	err := json.NewDecoder(r.Body).Decode(&updatedCinemaHall)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.repository.UpdateCinemaHall(&updatedCinemaHall)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, updatedCinemaHall, http.StatusOK)
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
