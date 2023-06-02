package repository

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/movie/service"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type movie struct {
	Id          int
	Title       string
	Genre       string
	ReleaseDate string
	Duration    int
}

type HallRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *HallRepository {
	return &HallRepository{db: db}
}

//func (h *HallRepository) Halls() ([]service.Hall, error) {
//	rows, err := h.db.Query(`SELECT hall_id, hall_name, capacity FROM halls`)
//	if err != nil {
//		log.Println(err)
//		return nil, err
//	}
//
//	defer func() {
//		if err = rows.Close(); err != nil {
//			log.Println(err)
//		}
//	}()
//
//	var cinemaHalls []service.Hall
//	for rows.Next() {
//		var hall hall
//		if err = rows.Scan(&hall.Id, &hall.Name, &hall.Capacity); err != nil {
//			log.Println(err)
//			return nil, fmt.Errorf("failed to get hall: %w", err)
//		}
//		cinemaHalls = append(cinemaHalls, service.NewHallEntity(hall.Id, hall.Name, hall.Capacity))
//	}
//
//	return cinemaHalls, nil
//}

func (h *HallRepository) MovieById(id int) (service.Movie, error) {
	row := h.db.QueryRow(`SELECT movie_id, title, genre, release_date, duration 
								FROM movies 
								WHERE movie_id = $1`, id)
	var movie movie
	err := row.Scan(&movie.Id, &movie.Title, &movie.Genre, &movie.ReleaseDate, &movie.ReleaseDate, &movie.Duration)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(err)
			return service.Movie{}, service.ErrMovieNotFound
		}
		log.Println(err)
		return service.Movie{}, fmt.Errorf("could not get movie by id: %w", err)
	}

	return service.NewHallEntity(hall.Id, hall.Name, hall.Capacity), nil
}

func (h *HallRepository) CreateMovie(title, genre, releaseDate string, duration int) (hallId int, err error) {
	var id int
	err = h.db.QueryRow(`INSERT INTO movies (title, genre, release_date, duration)
								VALUES ($1, $2, $3, $4)
								RETURNING movie_id`, title, genre, releaseDate, duration).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil
}

func (h *HallRepository) UpdateMovie(id int, title, genre, releaseDate string, duration int) error {
	_, err := h.db.Exec(`UPDATE movies
								SET title = $1, genre = $2, release_date = $3, duration = $4
								WHERE movie_id = $5`, title, genre, releaseDate, duration, id)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to update movie: %w", err)
	}

	return nil
}

func (h *HallRepository) DeleteMovie(id int) error {
	_, err := h.db.Exec(`DELETE FROM movies WHERE movie_id = $1`, id)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	return nil
}

func (h *HallRepository) MovieExists(id int) (bool, error) {
	var count int
	err := h.db.QueryRow(`SELECT COUNT(*) FROM movies WHERE movie_id = $1`, id).Scan(&count)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("failed to check if movie exists %w", err)
	}

	return count > 0, nil
}
