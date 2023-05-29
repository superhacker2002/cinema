package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/hall/handler"
)

var (
	ErrHallNotFound = errors.New("Cinema hall not found")
)

type HallRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *HallRepository {
	return &HallRepository{db: db}
}

func (r *HallRepository) GetCinemaHalls() ([]handler.CinemaHall, error) {
	rows, err := r.db.Query("SELECT hall_id, hall_name, capacity FROM halls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cinemaHalls []handler.CinemaHall
	for rows.Next() {
		var hall handler.CinemaHall
		err := rows.Scan(&hall.ID, &hall.Name, &hall.Capacity)
		if err != nil {
			return nil, err
		}
		cinemaHalls = append(cinemaHalls, hall)
	}

	return cinemaHalls, nil
}

func (r *HallRepository) GetCinemaHallByID(id int) (*handler.CinemaHall, error) {
	row := r.db.QueryRow("SELECT hall_id, hall_name, capacity, available FROM halls WHERE hall_id = $1", id)
	var hall handler.CinemaHall
	err := row.Scan(&hall.ID, &hall.Name, &hall.Capacity, &hall.Available)
	if err != nil {
		return nil, err
	}

	return &hall, nil
}

func (r *HallRepository) CreateCinemaHall(hall *handler.CinemaHall) error {
	query := "INSERT INTO halls (hall_name, capacity, available) VALUES ($1, $2, $3) RETURNING id"

	err := r.db.QueryRow(query, hall.Name, hall.Capacity, hall.Available).Scan(&hall.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *HallRepository) UpdateCinemaHall(hall *handler.CinemaHall) error {
	_, err := r.db.Exec("UPDATE halls SET hall_name = $1, capacity = $2 WHERE hall_id = $3",
		hall.Name, hall.Capacity, hall.ID)
	return err
}

func (r *HallRepository) DeleteCinemaHall(id int) error {
	_, err := r.db.Exec("DELETE FROM halls WHERE hall_id = $1", id)
	return err
}

func (r *HallRepository) UpdateHallAvailability(id int, available bool) error {
	_, err := r.db.Exec("UPDATE halls SET available = $1 WHERE hall_id = $2", available, id)
	return err
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

	w.WriteHeader(http.StatusOK)
}
