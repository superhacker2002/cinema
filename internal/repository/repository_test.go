package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
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
		t.Fatalf("Error getting movie: %v", fmt.Errorf("could not get movie: %w", err))
	}

	assert.Equal(t, expectedMovie, movie, "Expected movie does not match the actual movie")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err, "Unfulfilled expectation")
}
