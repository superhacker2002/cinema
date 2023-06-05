package service

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasession/entity"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockRepo struct {
	sessions      []entity.CinemaSession
	sessionExists bool
	movieExists   bool
	hallExists    bool
	hallBusy      bool
	seats         []int
	id            int
	err           error
}

func (m *mockRepo) AvailableSeats(sessionId int) ([]int, error) {
	return m.seats, m.err
}

func (m *mockRepo) UpdateSession(id, movieId, hallId int, startTime, endTime string, price float32) error {
	return m.err
}

func (m *mockRepo) SessionEndTime(id int, startTime string) (string, error) {
	return "", nil
}

func (m *mockRepo) HallIsBusy(sessionId, hallId int, startTime, endTime string) (bool, error) {
	return m.hallBusy, nil
}

func (m *mockRepo) CreateSession(movieId, hallId int, startTime, endTime string, price float32) (int, error) {
	return m.id, m.err
}

func (m *mockRepo) DeleteSession(id int) (bool, error) {
	return m.sessionExists, m.err
}

func (m *mockRepo) SessionExists(id int) (bool, error) {
	return m.sessionExists, nil
}

func (m *mockRepo) HallExists(id int) (bool, error) {
	return m.hallExists, nil
}

func (m *mockRepo) MovieExists(id int) (bool, error) {
	return m.movieExists, nil
}

func (m *mockRepo) AllSessions(date string, offset, limit int) ([]entity.CinemaSession, error) {
	return m.sessions, m.err
}

func (m *mockRepo) SessionsForHall(hallId int, date string) ([]entity.CinemaSession, error) {
	return m.sessions, m.err
}

func TestAllSessions(t *testing.T) {
	repo := mockRepo{}
	t.Run("successful sessions get", func(t *testing.T) {
		sessions := []entity.CinemaSession{{}, {}}
		repo.sessions = sessions
		repo.err = nil

		s := New(&repo)
		serviceSessions, err := s.AllSessions("2024-05-18", 0, 10)
		assert.NoError(t, err)
		assert.Len(t, serviceSessions, 2)
	})

	t.Run("no cinema sessions were found", func(t *testing.T) {
		var sessions []entity.CinemaSession
		repo.sessions = sessions
		repo.err = ErrCinemaSessionsNotFound

		s := New(&repo)
		serviceSessions, err := s.AllSessions("2024-05-18", 0, 10)
		assert.Equal(t, err, ErrCinemaSessionsNotFound)
		assert.Len(t, serviceSessions, 0)
	})

	t.Run("repository error", func(t *testing.T) {
		var sessions []entity.CinemaSession
		repo.sessions = sessions
		repo.err = errors.New("something went wrong")

		s := New(&repo)
		serviceSessions, err := s.AllSessions("2024-05-18", 0, 10)
		assert.Equal(t, ErrInternalError, err)
		assert.Len(t, serviceSessions, 0)
	})
}

func TestSessionsForHall(t *testing.T) {
	repo := mockRepo{}
	t.Run("successful sessions get", func(t *testing.T) {
		sessions := []entity.CinemaSession{{}, {}}
		repo.sessions = sessions
		repo.err = nil

		s := New(&repo)
		serviceSessions, err := s.AllSessions("2024-05-18", 0, 10)
		assert.NoError(t, err)
		assert.Len(t, serviceSessions, 2)
	})

	t.Run("no cinema sessions were found", func(t *testing.T) {
		var sessions []entity.CinemaSession
		repo.sessions = sessions
		repo.err = ErrCinemaSessionsNotFound

		s := New(&repo)
		serviceSessions, err := s.AllSessions("2024-05-18", 0, 10)
		assert.Equal(t, err, ErrCinemaSessionsNotFound)
		assert.Len(t, serviceSessions, 0)
	})

	t.Run("repository error", func(t *testing.T) {
		var sessions []entity.CinemaSession
		repo.sessions = sessions
		repo.err = errors.New("something went wrong")

		s := New(&repo)
		serviceSessions, err := s.AllSessions("2024-05-18", 0, 10)
		assert.Equal(t, ErrInternalError, err)
		assert.Len(t, serviceSessions, 0)
	})
}

func TestCreateSession(t *testing.T) {
	repo := mockRepo{}
	t.Run("successful session creation", func(t *testing.T) {
		repo.hallExists = true
		repo.movieExists = true
		repo.hallBusy = false
		repo.err = nil
		repo.id = 1

		s := New(&repo)
		id, err := s.CreateSession(1, 1, "2023-05-30 20:00:00 +04", 10.0)
		assert.NoError(t, err)
		assert.NotZero(t, id)
	})

	t.Run("hall does not exist", func(t *testing.T) {
		repo.hallExists = false

		s := New(&repo)
		id, err := s.CreateSession(1, 1, "2023-05-30 20:00:00 +04", 10.0)
		assert.ErrorIs(t, err, ErrHallNotFound)
		assert.Zero(t, id)
	})

	t.Run("movie does not exist", func(t *testing.T) {
		repo.hallExists = true
		repo.movieExists = false

		s := New(&repo)
		id, err := s.CreateSession(1, 1, "2023-05-30 20:00:00 +04", 10.0)
		assert.ErrorIs(t, err, ErrMovieNotFound)
		assert.Zero(t, id)
	})

	t.Run("hall is busy", func(t *testing.T) {
		repo.hallExists = true
		repo.movieExists = true
		repo.hallBusy = true

		s := New(&repo)
		id, err := s.CreateSession(1, 1, "2023-05-30 20:00:00 +04", 10.0)
		assert.ErrorIs(t, err, ErrHallIsBusy)
		assert.Zero(t, id)
	})

	t.Run("repository error", func(t *testing.T) {
		repo.hallExists = true
		repo.movieExists = true
		repo.hallBusy = false
		repo.err = errors.New("something went wrong")

		s := New(&repo)
		id, err := s.CreateSession(1, 1, "2023-05-30 20:00:00 +04", 10.0)
		assert.ErrorIs(t, err, ErrInternalError)
		assert.Zero(t, id)
	})
}

func TestDeleteSession(t *testing.T) {
	repo := mockRepo{}
	t.Run("successful session deletion", func(t *testing.T) {
		repo.sessionExists = true

		s := New(&repo)
		err := s.DeleteSession(1)
		assert.NoError(t, err)
	})

	t.Run("session does not exist", func(t *testing.T) {
		repo.sessionExists = false

		s := New(&repo)
		err := s.DeleteSession(1)
		assert.ErrorIs(t, err, ErrCinemaSessionsNotFound)
	})

	t.Run("repository error", func(t *testing.T) {
		repo.sessionExists = true
		repo.err = errors.New("something went wrong")

		s := New(&repo)
		err := s.DeleteSession(1)
		assert.ErrorIs(t, err, ErrInternalError)
	})
}

func TestUpdateSession(t *testing.T) {
	repo := mockRepo{}
	t.Run("successful session creation", func(t *testing.T) {
		repo.sessionExists = true
		repo.hallExists = true
		repo.movieExists = true
		repo.hallBusy = false
		repo.err = nil
		repo.id = 1

		s := New(&repo)
		err := s.UpdateSession(1, 1, 1, "2023-05-30 20:00:00 +04", 10.0)
		assert.NoError(t, err)
	})

	t.Run("session does not exist", func(t *testing.T) {
		repo.sessionExists = false

		s := New(&repo)
		err := s.UpdateSession(1, 1, 1, "2023-05-30 20:00:00 +04", 10.0)
		assert.ErrorIs(t, err, ErrCinemaSessionsNotFound)
	})

	t.Run("hall does not exist", func(t *testing.T) {
		repo.sessionExists = true
		repo.hallExists = false

		s := New(&repo)
		err := s.UpdateSession(1, 1, 1, "2023-05-30 20:00:00 +04", 10.0)
		assert.ErrorIs(t, err, ErrHallNotFound)
	})

	t.Run("movie does not exist", func(t *testing.T) {
		repo.sessionExists = true
		repo.hallExists = true
		repo.movieExists = false

		s := New(&repo)
		err := s.UpdateSession(1, 1, 1, "2023-05-30 20:00:00 +04", 10.0)
		assert.ErrorIs(t, err, ErrMovieNotFound)
	})

	t.Run("hall is busy", func(t *testing.T) {
		repo.sessionExists = true
		repo.hallExists = true
		repo.movieExists = true
		repo.hallBusy = true

		s := New(&repo)
		err := s.UpdateSession(1, 1, 1, "2023-05-30 20:00:00 +04", 10.0)
		assert.ErrorIs(t, err, ErrHallIsBusy)
	})

	t.Run("repository error", func(t *testing.T) {
		repo.sessionExists = true
		repo.hallExists = true
		repo.movieExists = true
		repo.hallBusy = false
		repo.err = errors.New("something went wrong")

		s := New(&repo)
		err := s.UpdateSession(1, 1, 1, "2023-05-30 20:00:00 +04", 10.0)
		assert.ErrorIs(t, err, ErrInternalError)
	})
}

func TestAvailableSeats(t *testing.T) {
	repo := mockRepo{}

	t.Run("successful seats get", func(t *testing.T) {
		repo.sessionExists = true
		repo.seats = []int{3, 2, 1}
		repo.err = nil

		s := New(&repo)
		seats, err := s.AvailableSeats(1)
		assert.NoError(t, err)
		assert.Equal(t, seats, []int{1, 2, 3})
	})

	t.Run("session does not exist", func(t *testing.T) {
		repo.sessionExists = false

		s := New(&repo)
		seats, err := s.AvailableSeats(1)
		assert.ErrorIs(t, err, ErrCinemaSessionsNotFound)
		assert.Zero(t, len(seats))
	})

	t.Run("no available seats", func(t *testing.T) {
		repo.sessionExists = true
		repo.err = ErrNoAvailableSeats

		s := New(&repo)
		seats, err := s.AvailableSeats(1)
		assert.ErrorIs(t, err, ErrNoAvailableSeats)
		assert.Zero(t, len(seats))
	})

	t.Run("repository error", func(t *testing.T) {
		repo.sessionExists = true
		repo.err = errors.New("something went wrong")

		s := New(&repo)
		seats, err := s.AvailableSeats(1)
		assert.ErrorIs(t, err, ErrInternalError)
		assert.Zero(t, len(seats))
	})
}
