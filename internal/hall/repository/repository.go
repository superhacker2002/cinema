package repository

import (
	"database/sql"
)

type HallRepository struct {
	db *sql.DB
}

func New(db *sql.DB) HallRepository {
	return HallRepository{db: db}
}

type Repository interface{}

type CinemaHall struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Capacity  int    `json:"capacity"`
	Available bool   `json:"available"`
}
