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

type MovieRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (h *MovieRepository) Movies(date string) ([]service.Movie, error) {
	rows, err := h.db.Query(`SELECT DISTINCT m.*
									FROM movies m
									JOIN cinema_sessions cs ON m.movie_id = cs.movie_id
									WHERE date_trunc('day', start_time) = $1`, date)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer func() {
		if err = rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	var movies []service.Movie
	for rows.Next() {
		var movie movie
		if err = rows.Scan(&movie.Id, &movie.Title, &movie.Genre, &movie.ReleaseDate, &movie.Duration); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("failed to get movie: %w", err)
		}
		movies = append(movies,
			service.NewMovieEntity(movie.Id, movie.Title, movie.Genre, movie.ReleaseDate, movie.Duration))
	}

	return movies, nil
}

func (h *MovieRepository) MovieById(id int) (service.Movie, error) {
	row := h.db.QueryRow(`SELECT *
								FROM movies 
								WHERE movie_id = $1`, id)
	var movie movie
	err := row.Scan(&movie.Id, &movie.Title, &movie.Genre, &movie.ReleaseDate, &movie.Duration)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(err)
			return service.Movie{}, service.ErrMovieNotFound
		}
		log.Println(err)
		return service.Movie{}, fmt.Errorf("could not get movie by id: %w", err)
	}

	return service.NewMovieEntity(movie.Id, movie.Title, movie.Genre, movie.ReleaseDate, movie.Duration), nil
}

func (h *MovieRepository) CreateMovie(title, genre, releaseDate string, duration int) (hallId int, err error) {
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

func (h *MovieRepository) UpdateMovie(id int, title, genre, releaseDate string, duration int) error {
	_, err := h.db.Exec(`UPDATE movies
								SET title = $1, genre = $2, release_date = $3, duration = $4
								WHERE movie_id = $5`, title, genre, releaseDate, duration, id)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to update movie: %w", err)
	}

	return nil
}

func (h *MovieRepository) DeleteMovie(id int) error {
	_, err := h.db.Exec(`DELETE FROM movies WHERE movie_id = $1`, id)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	return nil
}

func (h *MovieRepository) MovieExists(id int) (bool, error) {
	var count int
	err := h.db.QueryRow(`SELECT COUNT(*) FROM movies WHERE movie_id = $1`, id).Scan(&count)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("failed to check if movie exists %w", err)
	}

	return count > 0, nil
}
