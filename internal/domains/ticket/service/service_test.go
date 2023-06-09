package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mockRepository struct {
	sessionExists bool
	ticketExists  bool
	err           error
}

func (m *mockRepository) TicketExists(sessionId, seatNum int) (bool, error) {
	return m.ticketExists, nil
}

func (m *mockRepository) SessionExists(id int) (bool, error) {
	return m.sessionExists, nil
}

func (m *mockRepository) CreateTicket(sessionId, userId, seatNum int) (Ticket, error) {
	if sessionId == 1 && userId == 1 && seatNum == 2 {
		return NewTicketEntity(1, 1, 2, 120, "Movie 1", time.Now()), nil
	}
	if sessionId == 2 && userId == 2 && seatNum == 3 {
		return NewTicketEntity(2, 2, 3, 90, "Movie 2", time.Now()), nil
	}
	return Ticket{}, m.err
}

type mockTicketGenerator struct{}

func (m *mockTicketGenerator) GenerateTicket(t Ticket, outputPath string) error {
	if outputPath == "ticket1.pdf" {
		return nil
	}

	return nil
}

func TestService_BuyTicket(t *testing.T) {
	repo := &mockRepository{}
	gen := &mockTicketGenerator{}

	t.Run("successful purchase", func(t *testing.T) {
		repo.sessionExists = true
		service := New(repo, gen)
		outputPath, err := service.BuyTicket(1, 1, 2)
		assert.NoError(t, err)
		assert.Equal(t, "ticket1.pdf", outputPath)
	})

	t.Run("session not found", func(t *testing.T) {
		repo.sessionExists = false
		service := New(repo, gen)
		_, err := service.BuyTicket(3, 1, 2)
		assert.ErrorIs(t, err, ErrCinemaSessionsNotFound)
	})

	t.Run("ticket already exists", func(t *testing.T) {
		repo.ticketExists = true
		repo.sessionExists = true
		service := New(repo, gen)
		_, err := service.BuyTicket(4, 5, 6)
		assert.ErrorIs(t, err, ErrTicketExists)
	})

	t.Run("internal server error", func(t *testing.T) {
		repo.ticketExists = false
		repo.sessionExists = true
		repo.err = errors.New("something went wrong")
		service := New(repo, gen)
		_, err := service.BuyTicket(4, 4, 4)
		assert.ErrorIs(t, err, ErrInternalError)
	})
}
