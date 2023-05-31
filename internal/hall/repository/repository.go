package repository

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/hall/service"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type hall struct {
	Id       int
	Name     string
	Capacity int
}

type HallRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *HallRepository {
	return &HallRepository{db: db}
}

func (h *HallRepository) Halls() ([]service.Hall, error) {
	rows, err := h.db.Query(`SELECT hall_id, hall_name, capacity FROM halls`)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	var cinemaHalls []service.Hall
	for rows.Next() {
		var hall hall
		if err = rows.Scan(&hall.Id, &hall.Name, &hall.Capacity); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("failed to get hall: %w", err)
		}
		cinemaHalls = append(cinemaHalls, service.NewHallEntity(hall.Id, hall.Name, hall.Capacity))
	}

	return cinemaHalls, nil
}

func (h *HallRepository) HallById(id int) (service.Hall, error) {
	row := h.db.QueryRow(`SELECT hall_id, hall_name, capacity 
								FROM halls 
								WHERE hall_id = $1`, id)
	var hall hall
	err := row.Scan(&hall.Id, &hall.Name, &hall.Capacity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return service.Hall{}, service.ErrHallNotFound
		}
		return service.Hall{}, fmt.Errorf("could not get user credentials: %w", err)
	}

	return service.NewHallEntity(hall.Id, hall.Name, hall.Capacity), nil
}

func (h *HallRepository) CreateHall(name string, capacity int) (hallId int, err error) {
	var id int
	err = h.db.QueryRow(`INSERT INTO halls (hall_name, capacity)
								VALUES ($1, $2)
								RETURNING hall_id`, name, capacity).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil
}

func (h *HallRepository) UpdateCinemaHall(id int, name string, capacity int) error {
	_, err := h.db.Exec(`UPDATE halls
								SET hall_name = $1, capacity = $2
								WHERE hall_id = $3`, name, capacity, id)
	if err != nil {
		return fmt.Errorf("failed to update hall: %w", err)
	}

	return nil
}

func (h *HallRepository) DeleteCinemaHall(id int) error {
	_, err := h.db.Exec(`DELETE FROM halls WHERE hall_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete hall: %w", err)
	}

	return nil
}

func (h *HallRepository) HallExists(id int) (bool, error) {
	var count int
	err := h.db.QueryRow(`SELECT COUNT(*) FROM halls WHERE hall_id = $1`, id).Scan(&count)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("failed to check if hall exists %w", err)
	}

	return count > 0, nil
}

/*
func (h *HallRepository) UpdateHallAvailability(id int, available bool) error {
	_, err := h.db.Exec("UPDATE halls SET available = $1 WHERE hall_id = $2", available, id)
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
*/
