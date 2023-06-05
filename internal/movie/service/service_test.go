package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	movies      []Movie
	movieExists bool
	userExists  bool
	id          int
	err         error
}

func (m *mockRepository) WatchedMovies(userId int) (bool, error, []Movie) {
	return m.userExists, m.err, m.movies
}

func (m *mockRepository) Movies(_ string) ([]Movie, error) {
	return m.movies, m.err
}

func (m *mockRepository) MovieById(id int) (Movie, error) {
	switch id {
	case 1:
		return Movie{Id: 1,
			Title:       "Avengers: Endgame",
			Genre:       "Action, Adventure, Drama",
			ReleaseDate: "2019-04-26",
			Duration:    181,
		}, nil
	case 2:
		return Movie{Id: 1,
			Title:       "The Lion Kin",
			Genre:       "Animation, Adventure, Drama",
			ReleaseDate: "2019-07-19",
			Duration:    118,
		}, nil
	default:
		return Movie{}, ErrMoviesNotFound
	}
}

func (m *mockRepository) CreateMovie(title, genre, releaseDate string, duration int) (movieId int, err error) {
	return m.id, m.err
}

func (m *mockRepository) UpdateMovie(id int, title, genre, releaseDate string, duration int) (bool, error) {
	return m.movieExists, m.err
}

func (m *mockRepository) DeleteMovie(id int) (bool, error) {
	return m.movieExists, m.err
}

func TestMovies(t *testing.T) {
	repo := &mockRepository{}

	t.Run("successful movies get", func(t *testing.T) {
		movies := []Movie{
			{
				Id:          1,
				Title:       "Avengers: Endgame",
				Genre:       "Action, Adventure, Drama",
				ReleaseDate: "2019-04-26",
				Duration:    181,
			},
			{
				Id:          1,
				Title:       "The Lion Kin",
				Genre:       "Animation, Adventure, Drama",
				ReleaseDate: "2019-07-19",
				Duration:    118,
			},
		}
		repo.movies = movies
		s := New(repo)
		respMovies, err := s.Movies()
		assert.NoError(t, err)
		assert.Equal(t, movies, respMovies)
	})

	t.Run("repository error", func(t *testing.T) {
		repo.err = errors.New("something went wrong")
		s := New(repo)
		respMovies, err := s.Movies()
		assert.Error(t, ErrInternalError, err)
		assert.Zero(t, len(respMovies))
	})
}

func TestMovieById(t *testing.T) {
	repo := mockRepository{}
	t.Run("successful hall get", func(t *testing.T) {
		movie := Movie{
			Id:          1,
			Title:       "Avengers: Endgame",
			Genre:       "Action, Adventure, Drama",
			ReleaseDate: "2019-04-26",
			Duration:    181,
		}

		s := New(&repo)
		respMovie, err := s.MovieById(1)
		assert.NoError(t, err)
		assert.Equal(t, movie, respMovie)
	})

	t.Run("movie does not exist", func(t *testing.T) {
		s := New(&repo)
		_, err := s.MovieById(3)
		assert.ErrorIs(t, err, ErrMoviesNotFound)
	})
}

func TestCreateMovie(t *testing.T) {
	repo := mockRepository{}
	t.Run("successful movie creation", func(t *testing.T) {
		repo.id = 3
		s := New(&repo)
		id, err := s.CreateMovie("Movie", "Art house", "2023-05-30", 190)
		assert.NoError(t, err)
		assert.Equal(t, 3, id)
	})

	t.Run("repository error", func(t *testing.T) {
		repo.err = errors.New("something went wrong")
		s := New(&repo)
		id, err := s.CreateMovie("Movie", "Art house", "2023-05-30", 190)
		assert.ErrorIs(t, err, ErrInternalError)
		assert.Zero(t, id)
	})
}

func TestUpdateMovie(t *testing.T) {
	repo := mockRepository{}
	t.Run("successful movie update", func(t *testing.T) {
		repo.movieExists = true
		s := New(&repo)
		err := s.UpdateMovie(1, "Movie", "Art house", "2023-05-30", 190)
		assert.NoError(t, err)
	})

	t.Run("movie does not exist", func(t *testing.T) {
		repo.movieExists = false
		s := New(&repo)
		err := s.UpdateMovie(1, "Movie", "Art house", "2023-05-30", 190)
		assert.ErrorIs(t, err, ErrMoviesNotFound)
	})
}

func TestDeleteMovie(t *testing.T) {
	repo := mockRepository{}
	t.Run("successful movie delete", func(t *testing.T) {
		repo.movieExists = true
		s := New(&repo)
		err := s.DeleteMovie(1)
		assert.NoError(t, err)
	})

	t.Run("movie does not exist", func(t *testing.T) {
		repo.movieExists = false
		s := New(&repo)
		err := s.DeleteMovie(3)
		assert.ErrorIs(t, err, ErrMoviesNotFound)
	})
}

func TestWatchedMovies(t *testing.T) {
	repo := &mockRepository{}

	t.Run("successful movies get", func(t *testing.T) {
		movies := []Movie{
			{
				Id:          1,
				Title:       "Avengers: Endgame",
				Genre:       "Action, Adventure, Drama",
				ReleaseDate: "2019-04-26",
				Duration:    181,
			},
			{
				Id:          1,
				Title:       "The Lion Kin",
				Genre:       "Animation, Adventure, Drama",
				ReleaseDate: "2019-07-19",
				Duration:    118,
			},
		}
		repo.movies = movies
		repo.userExists = true
		s := New(repo)
		respMovies, err := s.WatchedMovies(1)
		assert.NoError(t, err)
		assert.Equal(t, movies, respMovies)
	})

	t.Run("repository error", func(t *testing.T) {
		repo.err = errors.New("something went wrong")
		repo.userExists = true
		s := New(repo)
		respMovies, err := s.WatchedMovies(1)
		assert.Error(t, ErrInternalError, err)
		assert.Zero(t, len(respMovies))
	})

	t.Run("user does not exist", func(t *testing.T) {
		repo.userExists = false
		s := New(repo)
		respMovies, err := s.WatchedMovies(1)
		assert.Error(t, ErrUserNotFound, err)
		assert.Zero(t, len(respMovies))
	})
}
