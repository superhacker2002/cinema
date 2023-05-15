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
