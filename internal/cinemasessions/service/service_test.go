package service

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/entity"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockRepo struct {
	sessions []entity.CinemaSession
	err      error
}

const layout = "2006-01-02 15:04:05 MST"

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
