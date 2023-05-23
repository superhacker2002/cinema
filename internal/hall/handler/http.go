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
	db         *sql.DB
}

func NewHTTPHandler(router *mux.Router, repository repository.Repository, db *sql.DB) {
	handler := HTTPHandler{repository: repository, db: db}
	handler.setRoutes(router)
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
	rows, err := h.db.Query("SELECT hall_id, hall_name, capacity FROM halls")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cinemaHalls []repository.CinemaHall
	for rows.Next() {
		var hall repository.CinemaHall
		err := rows.Scan(&hall.ID, &hall.Name, &hall.Capacity)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		cinemaHalls = append(cinemaHalls, hall)
	}

	api.WriteResponse(w, cinemaHalls, http.StatusOK)
}

func (h HTTPHandler) createHallHandler(w http.ResponseWriter, r *http.Request) {
	var newCinemaHall repository.CinemaHall
	err := json.NewDecoder(r.Body).Decode(&newCinemaHall)
	if err != nil {
		api.HandleError(w, err, http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec("INSERT INTO halls (hall_name, capacity, available) VALUES ($1, $2, $3)",
		newCinemaHall.Name, newCinemaHall.Capacity, newCinemaHall.Available)
	if err != nil {
		api.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		api.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	newCinemaHall.ID = int(rowsAffected)
	api.WriteResponse(w, newCinemaHall, http.StatusOK)
}

func (h HTTPHandler) getHallHandler(w http.ResponseWriter, r *http.Request) {
	hallID, err := api.GetHallID(r)
	if err != nil {
		api.HandleError(w, err, http.StatusBadRequest)
		return
	}

	row := h.db.QueryRow("SELECT hall_id, hall_name, capacity, available FROM halls WHERE hall_id = $1", hallID)
	var hall repository.CinemaHall
	err = row.Scan(&hall.ID, &hall.Name, &hall.Capacity, &hall.Available)
	if err != nil {
		if err == sql.ErrNoRows {
			api.HandleError(w, fmt.Errorf("Hall notfound"), http.StatusNotFound)
		} else {
			api.HandleError(w, err, http.StatusInternalServerError)
		}
		return
	}

	api.WriteResponse(w, hall, http.StatusOK)
}

func (h HTTPHandler) updateHallHandler(w http.ResponseWriter, r *http.Request) {
	var updatedCinemaHall repository.CinemaHall
	err := json.NewDecoder(r.Body).Decode(&updatedCinemaHall)
	if err != nil {
		api.HandleError(w, err, http.StatusBadRequest)
		return
	}

	_, err = h.db.Exec("UPDATE halls SET hall_name = $1, capacity = $2 WHERE hall_id = $3",
		updatedCinemaHall.Name, updatedCinemaHall.Capacity, updatedCinemaHall.ID)
	if err != nil {
		api.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	api.WriteResponse(w, updatedCinemaHall, http.StatusOK)
}

func (h HTTPHandler) deleteHallHandler(w http.ResponseWriter, r *http.Request) {
	hallID, err := api.GetHallID(r)
	if err != nil {
		api.HandleError(w, err, http.StatusBadRequest)
		return
	}

	_, err = h.db.Exec("DELETE FROM halls WHERE hall_id = $1", hallID)
	if err != nil {
		api.HandleError(w, err, http.StatusInternalServerError)
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
		api.HandleError(w, err, http.StatusBadRequest)
		return
	}

	_, err = h.db.Exec("UPDATE halls SET available = $1 WHERE hall_id = $2",
		update.Available, update.HallID)
	if err != nil {
		api.HandleError(w, err, http.StatusInternalServerError)
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
		api.HandleError(w, err, http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE halls SET assigned_movie = $1, seats_available = $2 WHERE hall_id = $3",
		assignment.Movie, assignment.Seats, assignment.HallID)
	if err != nil {
		api.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Assigned movie '%s' to hall with ID %d", assignment.Movie, assignment.HallID)
	api.WriteResponse(w, message, http.StatusOK)
}
