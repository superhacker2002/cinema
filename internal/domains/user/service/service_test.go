package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockRepository struct {
	id  int
	err error
}

func (m mockRepository) CreateUser(username string, passwordHash string) (userId int, err error) {
	return m.id, m.err
}

func TestCreateUser(t *testing.T) {
	repo := mockRepository{}
	t.Run("successful user creation", func(t *testing.T) {
		repo.id = 3
		s := New(&repo)
		id, err := s.CreateUser("test_user", "hashed_password")
		assert.NoError(t, err)
		assert.Equal(t, 3, id)
	})

	t.Run("repository error", func(t *testing.T) {
		repo.err = errors.New("something went wrong")
		s := New(&repo)
		id, err := s.CreateUser("test_user", "hashed_password")
		assert.ErrorIs(t, err, ErrInternalError)
		assert.Zero(t, id)
	})

	t.Run("user exists", func(t *testing.T) {
		repo.err = ErrUserExists
		s := New(&repo)
		id, err := s.CreateUser("test_user", "hashed_password")
		assert.ErrorIs(t, err, ErrUserExists)
		assert.Zero(t, id)
	})
}
