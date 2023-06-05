package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	halls      []Hall
	hallExists bool
	id         int
	err        error
}

func (m *mockRepository) Halls() ([]Hall, error) {
	return m.halls, m.err
}

func (m *mockRepository) HallById(id int) (Hall, error) {
	switch id {
	case 1:
		return Hall{Id: 1, Name: "Hall 1", Capacity: 100}, nil
	case 2:
		return Hall{Id: 2, Name: "Hall 2", Capacity: 150}, nil
	default:
		return Hall{}, ErrHallNotFound
	}
}

func (m *mockRepository) CreateHall(name string, capacity int) (int, error) {
	return m.id, m.err
}

func (m *mockRepository) UpdateHall(id int, name string, capacity int) (bool, error) {
	return m.hallExists, m.err
}

func (m *mockRepository) DeleteHall(id int) (bool, error) {
	return m.hallExists, m.err
}

func TestHalls(t *testing.T) {
	repo := &mockRepository{}

	t.Run("successful halls get", func(t *testing.T) {
		halls := []Hall{
			{Id: 1, Name: "Hall 1", Capacity: 100},
			{Id: 2, Name: "Hall 2", Capacity: 150},
		}
		repo.halls = halls
		s := New(repo)
		respHalls, err := s.Halls()
		assert.NoError(t, err)
		assert.Equal(t, halls, respHalls)
	})

	t.Run("repository error", func(t *testing.T) {
		repo.err = errors.New("something went wrong")
		s := New(repo)
		respHalls, err := s.Halls()
		assert.Error(t, ErrInternalError, err)
		assert.Zero(t, len(respHalls))
	})
}

func TestHallById(t *testing.T) {
	repo := mockRepository{}
	t.Run("successful hall get", func(t *testing.T) {
		hall := Hall{Id: 1, Name: "Hall 1", Capacity: 100}

		s := New(&repo)
		respHall, err := s.HallById(1)
		assert.NoError(t, err)
		assert.Equal(t, hall, respHall)
	})

	t.Run("hall does not exist", func(t *testing.T) {
		s := New(&repo)
		_, err := s.HallById(3)
		assert.ErrorIs(t, err, ErrHallNotFound)
	})
}

func TestCreateHall(t *testing.T) {
	repo := mockRepository{}
	t.Run("successful hall creation", func(t *testing.T) {
		repo.id = 3
		s := New(&repo)
		id, err := s.CreateHall("Hall 3", 200)
		assert.NoError(t, err)
		assert.Equal(t, 3, id)
	})

	t.Run("repository error", func(t *testing.T) {
		repo.err = errors.New("something went wrong")
		s := New(&repo)
		id, err := s.CreateHall("Hall 3", 200)
		assert.ErrorIs(t, err, ErrInternalError)
		assert.Zero(t, id)
	})
}

func TestUpdateHall(t *testing.T) {
	repo := mockRepository{}
	t.Run("successful hall update", func(t *testing.T) {
		repo.hallExists = true
		s := New(&repo)
		err := s.UpdateHall(1, "Hall 3", 200)
		assert.NoError(t, err)
	})

	t.Run("hall does not exist", func(t *testing.T) {
		repo.hallExists = false
		s := New(&repo)
		err := s.UpdateHall(1, "Hall 3", 200)
		assert.ErrorIs(t, err, ErrHallNotFound)
	})
}

func TestDeleteHall(t *testing.T) {
	repo := mockRepository{}
	t.Run("successful hall delete", func(t *testing.T) {
		repo.hallExists = true
		s := New(&repo)
		err := s.DeleteHall(1)
		assert.NoError(t, err)
	})

	t.Run("hall does not exist", func(t *testing.T) {
		repo.hallExists = false
		s := New(&repo)
		err := s.DeleteHall(3)
		assert.ErrorIs(t, err, ErrHallNotFound)
	})

	t.Run("repository error", func(t *testing.T) {
		repo.hallExists = true
		repo.err = errors.New("something went wrong")
		s := New(&repo)
		err := s.DeleteHall(3)
		assert.ErrorIs(t, err, ErrInternalError)
	})
}
