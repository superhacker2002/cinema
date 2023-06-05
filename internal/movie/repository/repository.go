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

func (m *MovieRepository) Movies(date string) ([]service.Movie, error) {
	rows, err := m.db.Query(`SELECT DISTINCT m.*
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

	movies, err := m.readMovies(rows)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return movies, nil
}

func (m *MovieRepository) MovieById(id int) (service.Movie, error) {
	row := m.db.QueryRow(`SELECT *
						FROM movies 
						WHERE movie_id = $1`, id)
	var movie movie
	err := row.Scan(&movie.Id, &movie.Title, &movie.Genre, &movie.ReleaseDate, &movie.Duration)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(err)
			return service.Movie{}, service.ErrMoviesNotFound
		}
		log.Println(err)
		return service.Movie{}, fmt.Errorf("could not get movie by id: %w", err)
	}

	return service.NewMovieEntity(movie.Id, movie.Title, movie.Genre, movie.ReleaseDate, movie.Duration), nil
}

func (m *MovieRepository) CreateMovie(title, genre, releaseDate string, duration int) (mallId int, err error) {
	var id int
	err = m.db.QueryRow(`INSERT INTO movies (title, genre, release_date, duration)
						VALUES ($1, $2, $3, $4)
						RETURNING movie_id`, title, genre, releaseDate, duration).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil
}

func (m *MovieRepository) UpdateMovie(id int, title, genre, releaseDate string, duration int) (bool, error) {
	res, err := m.db.Exec(`UPDATE movies
						SET title = $1, genre = $2, release_date = $3, duration = $4
						WHERE movie_id = $5`, title, genre, releaseDate, duration, id)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("failed to update movie: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (m *MovieRepository) DeleteMovie(id int) (bool, error) {
	res, err := m.db.Exec(`DELETE FROM movies WHERE movie_id = $1`, id)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("failed to delete movie: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (m *MovieRepository) WatchedMovies(userId int) (bool, error, []service.Movie) {
	rows, err := m.db.Query(`SELECT DISTINCT m.*
							FROM movies m
							JOIN cinema_sessions cs ON m.movie_id = cs.movie_id
							JOIN tickets t ON cs.session_id = t.session_id
							WHERE t.user_id = $1;
							`, userId)

	if errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return false, nil, nil
	}

	if err != nil {
		log.Println(err)
		return false, err, nil
	}

	defer func() {
		if err = rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	movies, err := m.readMovies(rows)
	if err != nil {
		log.Println(err)
		return false, err, nil
	}

	return false, nil, movies
}

func (m *MovieRepository) readMovies(rows *sql.Rows) ([]service.Movie, error) {
	var movies []service.Movie
	for rows.Next() {
		var movie movie
		if err := rows.Scan(&movie.Id, &movie.Title, &movie.Genre, &movie.ReleaseDate, &movie.Duration); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("failed to get movie: %w", err)
		}
		movies = append(movies,
			service.NewMovieEntity(movie.Id, movie.Title, movie.Genre, movie.ReleaseDate, movie.Duration))
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error while iterating over movies: %w", err)
	}

	if len(movies) == 0 {
		log.Println(service.ErrMoviesNotFound)
		return nil, service.ErrMoviesNotFound
	}

	return movies, nil
}
