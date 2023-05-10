package repository

import (
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func TestGetMovie(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	movieID := 1
	expectedMovie := &Movie{
		ID:          movieID,
		Title:       "Test Movie",
		Genre:       "Action",
		ReleaseDate: pq.NullTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true},
		Duration:    120,
	}

	mock.ExpectQuery(`SELECT id, title, genre, release_date, duration FROM movies WHERE id = ?`).
		WithArgs(movieID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "genre", "release_date", "duration"}).
			AddRow(expectedMovie.ID, expectedMovie.Title, expectedMovie.Genre, expectedMovie.ReleaseDate, expectedMovie.Duration))

	repo := NewMovieRepository(db)

	movie, err := repo.GetMovie(movieID)
	if err != nil {
		t.Fatalf("Error getting movie: %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("Unfulfilled expectation: %v", err)
	}

	if !reflect.DeepEqual(expectedMovie, movie) {
		t.Errorf("Expected movie %v, but got %v", expectedMovie, movie)
	}
}
