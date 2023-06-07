package repository

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/hall/service"
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
		log.Println(err)
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
			log.Println(err)
			return service.Hall{}, service.ErrHallNotFound
		}
		log.Println(err)
		return service.Hall{}, fmt.Errorf("could not get hall by id: %w", err)
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

func (h *HallRepository) UpdateHall(id int, name string, capacity int) (bool, error) {
	res, err := h.db.Exec(`UPDATE halls
						SET hall_name = $1, capacity = $2
						WHERE hall_id = $3`, name, capacity, id)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("failed to update hall: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (h *HallRepository) DeleteHall(id int) (bool, error) {
	res, err := h.db.Exec(`DELETE FROM halls WHERE hall_id = $1`, id)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("failed to delete hall: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return false, nil
	}
	return true, nil
}
