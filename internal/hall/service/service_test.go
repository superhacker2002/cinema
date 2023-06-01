package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/hall/service"
)

type mockRepository struct{}

func (m *mockRepository) Halls() ([]service.Hall, error) {
	halls := []service.Hall{
		{Id: 1, Name: "Hall 1", Capacity: 100},
		{Id: 2, Name: "Hall 2", Capacity: 150},
	}
	return halls, nil
}

func (m *mockRepository) HallById(id int) (service.Hall, error) {
	switch id {
	case 1:
		return service.Hall{Id: 1, Name: "Hall 1", Capacity: 100}, nil
	case 2:
		return service.Hall{Id: 2, Name: "Hall 2", Capacity: 150}, nil
	default:
		return service.Hall{}, service.ErrHallNotFound
	}
}

func (m *mockRepository) CreateHall(name string, capacity int) (int, error) {
	// Implement create hall logic if needed
	return 0, nil
}

func (m *mockRepository) UpdateHall(id int, name string, capacity int) error {
	// Implement update hall logic if needed
	return nil
}

func (m *mockRepository) DeleteHall(id int) error {
	// Implement delete hall logic if needed
	return nil
}

func (m *mockRepository) HallExists(id int) (bool, error) {
	// Implement hall exists logic if needed
	return true, nil
}

func TestService(t *testing.T) {
	repo := &mockRepository{}
	svc := service.Service{R: repo}

	expectedHalls := []service.Hall{
		{Id: 1, Name: "Hall 1", Capacity: 100},
		{Id: 2, Name: "Hall 2", Capacity: 150},
	}

	// Test Halls
	halls, err := svc.Halls()
	assert.NoError(t, err)
	assert.Equal(t, expectedHalls, halls)

	// Test HallById (existing ID)
	expectedHall := service.Hall{Id: 1, Name: "Hall 1", Capacity: 100}
	hall, err := svc.HallById(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedHall, hall)

	// Test HallById (non-existing ID)
	_, err = svc.HallById(3)
	assert.ErrorIs(t, err, service.ErrHallNotFound)

	// Test CreateHall
	hallId, err := svc.CreateHall("New Hall", 200)
	assert.NoError(t, err)
	assert.NotZero(t, hallId)

	// Test UpdateHall (existing ID)
	err = svc.UpdateHall(1, "Updated Hall", 200)
	assert.NoError(t, err)

	// Test UpdateHall (non-existing ID)
	err = svc.UpdateHall(3, "Updated Hall", 200)
	assert.ErrorIs(t, err, service.ErrHallNotFound)

	// Test DeleteHall (existing ID)
	err = svc.DeleteHall(1)
	assert.NoError(t, err)

	// Test DeleteHall (non-existing ID)
	err = svc.DeleteHall(3)
	assert.ErrorIs(t, err, service.ErrHallNotFound)
}
