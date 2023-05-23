package repository

import (
	"database/sql"
	"errors"
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

func (r *HallRepository) GetCinemaHalls() ([]CinemaHall, error) {
	rows, err := r.db.Query("SELECT hall_id, hall_name, capacity FROM halls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cinemaHalls []CinemaHall
	for rows.Next() {
		var hall CinemaHall
		err := rows.Scan(&hall.ID, &hall.Name, &hall.Capacity)
		if err != nil {
			return nil, err
		}
		cinemaHalls = append(cinemaHalls, hall)
	}

	return cinemaHalls, nil
}

func (r *HallRepository) GetCinemaHallByID(id int) (*CinemaHall, error) {
	row := r.db.QueryRow("SELECT hall_id, hall_name, capacity, available FROM halls WHERE hall_id = $1", id)
	var hall CinemaHall
	err := row.Scan(&hall.ID, &hall.Name, &hall.Capacity, &hall.Available)
	if err != nil {
		return nil, err
	}

	return &hall, nil
}

func (r *HallRepository) CreateCinemaHall(hall *CinemaHall) error {
	result, err := r.db.Exec("INSERT INTO halls (hall_name, capacity, available) VALUES ($1, $2, $3)",
		hall.Name, hall.Capacity, hall.Available)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	hall.ID = int(rowsAffected)
	return nil
}

func (r *HallRepository) UpdateCinemaHall(hall *CinemaHall) error {
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
