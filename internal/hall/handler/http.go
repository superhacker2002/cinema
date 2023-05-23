package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/api"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/hall/repository"
	"github.com/gorilla/mux"
)

type HTTPHandler struct {
	repository repository.Repository
}

func New(router *mux.Router, repository repository.Repository) HTTPHandler {
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

func (h HTTPHandler) getHallsHandler(w http.ResponseWriter, r *http.Request) {
	cinemaHalls, err := h.repository.GetCinemaHalls()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.WriteResponse(w, cinemaHalls, http.StatusOK)
}

func (h HTTPHandler) createHallHandler(w http.ResponseWriter, r *http.Request) {
	var newCinemaHall repository.CinemaHall
	err := json.NewDecoder(r.Body).Decode(&newCinemaHall)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.repository.CreateCinemaHall(&newCinemaHall)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.WriteResponse(w, newCinemaHall, http.StatusOK)
}

func (h HTTPHandler) getHallHandler(w http.ResponseWriter, r *http.Request) {
	hallID, err := api.GetHallID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hall, err := h.repository.GetCinemaHallByID(hallID)
	if err != nil {
		if err == repository.ErrHallNotFound {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	api.WriteResponse(w, hall, http.StatusOK)
}

func (h HTTPHandler) updateHallHandler(w http.ResponseWriter, r *http.Request) {
	var updatedCinemaHall repository.CinemaHall
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

	api.WriteResponse(w, updatedCinemaHall, http.StatusOK)
}

func (h HTTPHandler) deleteHallHandler(w http.ResponseWriter, r *http.Request) {
	hallID, err := api.GetHallID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.repository.DeleteCinemaHall(hallID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Cinema hall with ID %d deleted", hallID)
	api.WriteResponse(w, message, http.StatusOK)
}

func (h HTTPHandler) updateAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	var update struct {
		HallID    int  `json:"hallId"`
		Available bool `json:"available"`
	}

	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.repository.UpdateHallAvailability(update.HallID, update.Available)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Updated availability for cinema hall with ID %d", update.HallID)
	api.WriteResponse(w, message, http.StatusOK)
}

func AssignMovie(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var assignment struct {
		HallID int    `json:"hallId"`
		Movie  string `json:"movie"`
		Seats  int    `json:"seats"`
	}

	err := json.NewDecoder(r.Body).Decode(&assignment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("UPDATE halls SET assigned_movie = $1, seats_available = $2 WHERE hall_id = $3",
		assignment.Movie, assignment.Seats, assignment.HallID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Assigned movie '%s' to hall with ID %d", assignment.Movie, assignment.HallID)
	api.WriteResponse(w, message, http.StatusOK)
}
