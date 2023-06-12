package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
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
	return Ticket{}, m.err
}

type mockTicketGenerator struct{}

func (m *mockTicketGenerator) GenerateTicket(t Ticket, w io.Writer) error {
	return nil
}

type mockTicketsStorage struct{}

func (m mockTicketsStorage) Store(ctx context.Context, file *os.File) (string, error) {
	return "", nil
}

func TestService_BuyTicket(t *testing.T) {
	repo := &mockRepository{}
	gen := &mockTicketGenerator{}
	storage := &mockTicketsStorage{}
	ctx := context.Background()

	t.Run("successful purchase", func(t *testing.T) {
		repo.sessionExists = true
		service := New(repo, gen, storage)
		_, err := service.BuyTicket(ctx, 1, 1, 2)
		assert.NoError(t, err)
	})

	t.Run("session not found", func(t *testing.T) {
		repo.sessionExists = false
		service := New(repo, gen, storage)
		_, err := service.BuyTicket(ctx, 1, 1, 2)
		assert.ErrorIs(t, err, ErrCinemaSessionsNotFound)
	})

	t.Run("ticket already exists", func(t *testing.T) {
		repo.ticketExists = true
		repo.sessionExists = true
		service := New(repo, gen, storage)
		_, err := service.BuyTicket(ctx, 1, 1, 2)
		assert.ErrorIs(t, err, ErrTicketExists)
	})

	t.Run("internal server error", func(t *testing.T) {
		repo.ticketExists = false
		repo.sessionExists = true
		repo.err = errors.New("something went wrong")

		service := New(repo, gen, storage)
		_, err := service.BuyTicket(ctx, 2, 1, 2)
		assert.ErrorIs(t, err, ErrInternalError)
	})
}
