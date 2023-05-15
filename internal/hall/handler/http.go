package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type repository interface{}

type HTTPHandler struct {
	repository repository
}

var db *sql.DB

type CinemaHall struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Capacity  int    `json:"capacity"`
	Available bool   `json:"available"`
}

func NewHTTPHandler(router *mux.Router, repository repository) HTTPHandler {
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
	rows, err := db.Query("SELECT hall_id, hall_name, capacity FROM halls")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cinemaHalls []CinemaHall
	for rows.Next() {
		var hall CinemaHall
		err := rows.Scan(&hall.ID, &hall.Name, &hall.Capacity)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		cinemaHalls = append(cinemaHalls, hall)
	}

	writeResponse(w, cinemaHalls, http.StatusOK)
}

func (h HTTPHandler) createHallHandler(w http.ResponseWriter, r *http.Request) {
	var newCinemaHall CinemaHall
	err := json.NewDecoder(r.Body).Decode(&newCinemaHall)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO halls (hall_name, capacity, available) VALUES ($1, $2, $3)",
		newCinemaHall.Name, newCinemaHall.Capacity, newCinemaHall.Available)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	newCinemaHall.ID = int(rowsAffected)
	writeResponse(w, newCinemaHall, http.StatusOK)
}

func (h HTTPHandler) getHallHandler(w http.ResponseWriter, r *http.Request) {
	hallID, err := getHallID(r)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT hall_id, hall_name, capacity, available FROM halls WHERE hall_id = $1", hallID)
	var hall CinemaHall
	err = row.Scan(&hall.ID, &hall.Name, &hall.Capacity, &hall.Available)
	if err != nil {
		if err == sql.ErrNoRows {
			handleError(w, fmt.Errorf("Hall notfound"), http.StatusNotFound)
		} else {
			handleError(w, err, http.StatusInternalServerError)
		}
		return
	}

	writeResponse(w, hall, http.StatusOK)
}

func (h HTTPHandler) updateHallHandler(w http.ResponseWriter, r *http.Request) {
	var updatedCinemaHall CinemaHall
	err := json.NewDecoder(r.Body).Decode(&updatedCinemaHall)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE halls SET hall_name = $1, capacity = $2 WHERE hall_id = $3",
		updatedCinemaHall.Name, updatedCinemaHall.Capacity, updatedCinemaHall.ID)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	writeResponse(w, updatedCinemaHall, http.StatusOK)
}

func (h HTTPHandler) deleteHallHandler(w http.ResponseWriter, r *http.Request) {
	hallID, err := getHallID(r)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM halls WHERE hall_id = $1", hallID)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Cinema hall with ID %d deleted", hallID)
	writeResponse(w, message, http.StatusOK)
}

func (h HTTPHandler) updateAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	var update struct {
		HallID    int  `json:"hallId"`
		Available bool `json:"available"`
	}

	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE halls SET available = $1 WHERE hall_id = $2",
		update.Available, update.HallID)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Updated availability for cinema hall with ID %d", update.HallID)
	writeResponse(w, message, http.StatusOK)
}

func assignMovie(w http.ResponseWriter, r *http.Request) {
	var assignment struct {
		HallID int    `json:"hallId"`
		Movie  string `json:"movie"`
		Seats  int    `json:"seats"`
	}

	err := json.NewDecoder(r.Body).Decode(&assignment)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE halls SET assigned_movie = $1, seats_available = $2 WHERE hall_id = $3",
		assignment.Movie, assignment.Seats, assignment.HallID)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Assigned movie '%s' to hall with ID %d", assignment.Movie, assignment.HallID)
	writeResponse(w, message, http.StatusOK)
}

// Helper function to retrieve the hall ID from the request URL parameters
func getHallID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	hallID, err := strconv.Atoi(vars["hallID"])
	if err != nil {
		return 0, fmt.Errorf("Invalid hall ID")
	}
	return hallID, nil
}

// Helper function to write a JSON response with the specified status code
func writeResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Helper function to handle errors and write an appropriate response
func handleError(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, err.Error(), statusCode)
}
