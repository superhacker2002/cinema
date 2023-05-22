package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/repository"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockRepo struct {
	sessions []repository.CinemaSession
	hallId   int
	err      error
}

func (m *mockRepo) AllSessions(date string, offset, limit int) ([]repository.CinemaSession, error) {
	return m.sessions, m.err
}

func (m *mockRepo) SessionsForHall(hallId int, date string) ([]repository.CinemaSession, error) {
	if hallId != m.hallId {
		return nil, repository.ErrCinemaSessionsNotFound
	}
	return m.sessions, m.err
}

func TestGetSessionsHandler(t *testing.T) {
	repo := mockRepo{}
	t.Run("successful sessions get", func(t *testing.T) {
		session := []repository.CinemaSession{
			{
				ID:        1,
				MovieId:   1,
				StartTime: "2024-05-18 20:00:00",
				EndTime:   "2024-05-18 22:00:00",
				Status:    "scheduled",
			},
		}
		repo.sessions = session
		repo.hallId = 1
		repo.err = nil

		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/1", nil)
		req = mux.SetURLVars(req, map[string]string{"hallId": "1"})
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{r: &repo}.getSessionsHandler
		handler(response, req)

		assert.NotEmpty(t, response.Body.String())
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("invalid hall id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/0", nil)
		req = mux.SetURLVars(req, map[string]string{"hallId": "0"})
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{r: &repo}.getSessionsHandler
		handler(response, req)

		assert.Equal(t, fmt.Sprintf("%v: 0\n", ErrInvalidHallId), response.Body.String())
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("no cinema sessions", func(t *testing.T) {
		repo.sessions = []repository.CinemaSession{{}}
		repo.hallId = 1
		repo.err = repository.ErrCinemaSessionsNotFound

		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/2", nil)
		req = mux.SetURLVars(req, map[string]string{"hallId": "2"})
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{r: &repo}.getSessionsHandler
		handler(response, req)

		assert.Equal(t, fmt.Sprintf("%v for hall 2\n", repository.ErrCinemaSessionsNotFound), response.Body.String())
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("repository error", func(t *testing.T) {
		repo.sessions = []repository.CinemaSession{{}}
		repo.hallId = 2
		repo.err = errors.New("something went wrong")

		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/2", nil)
		req = mux.SetURLVars(req, map[string]string{"hallId": "2"})
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{r: &repo}.getSessionsHandler
		handler(response, req)

		assert.Equal(t, fmt.Sprintf("%v\n", ErrInternalError), response.Body.String())
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestGetAllSessionsHandler(t *testing.T) {
	repo := mockRepo{}
	t.Run("successful sessions get", func(t *testing.T) {
		sessions := []repository.CinemaSession{
			{
				ID:        1,
				MovieId:   1,
				StartTime: "2024-05-18 20:00:00",
				EndTime:   "2024-05-18 22:00:00",
				Status:    "scheduled",
			},
		}
		repo.sessions = sessions
		repo.err = nil

		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/", nil)
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{r: &repo}.getAllSessionsHandler
		handler(response, req)

		assert.NotEmpty(t, response.Body.String())
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("no cinema sessions", func(t *testing.T) {
		sessions := []repository.CinemaSession{{}}
		repo.sessions = sessions
		repo.err = repository.ErrCinemaSessionsNotFound

		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/", nil)
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{r: &repo}.getAllSessionsHandler
		handler(response, req)

		assert.Equal(t, fmt.Sprintf("%v for all halls\n", repository.ErrCinemaSessionsNotFound), response.Body.String())
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("repository error", func(t *testing.T) {
		sessions := []repository.CinemaSession{{}}
		repo.sessions = sessions
		repo.err = errors.New("something went wrong")

		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/", nil)
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{r: &repo}.getAllSessionsHandler
		handler(response, req)

		assert.Equal(t, fmt.Sprintf("%v\n", ErrInternalError), response.Body.String())
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestPage(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "resource?offset=10&limit=20", nil)
		require.NoError(t, err, "failed to create test request")
		page, err := page(req)
		assert.NoError(t, err)
		assert.Equal(t, 10, page.Offset)
		assert.Equal(t, 20, page.Limit)
	})

	t.Run("invalid request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "resource?offset=-10&limit=20", nil)
		require.NoError(t, err, "failed to create test request")
		page, err := page(req)
		assert.Error(t, err)
		assert.Equal(t, 20, page.Limit)
	})

	t.Run("not a number offset", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "resource?offset=invalid&limit=20", nil)
		require.NoError(t, err, "failed to create test request")
		page, err := page(req)
		assert.Error(t, err)
		assert.Empty(t, page)
	})

	t.Run("empty offset", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "resource?limit=20", nil)
		require.NoError(t, err, "failed to create test request")
		page, err := page(req)
		assert.NoError(t, err)
		assert.Equal(t, 0, page.Offset)
		assert.Equal(t, 10, page.Limit)
	})
}

func TestDate(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "resource?date=2023-01-05", nil)
		require.NoError(t, err, "failed to create test request")
		date, err := date(req)
		assert.NoError(t, err)
		assert.Equal(t, date, "2023-01-05")
	})

	t.Run("invalid date", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "resource?date=2023/01/05", nil)
		require.NoError(t, err, "failed to create test request")
		date, err := date(req)
		assert.Error(t, err)
		assert.Equal(t, date, "")
	})

	t.Run("missing date", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "resource/", nil)
		require.NoError(t, err, "failed to create test request")
		date, err := date(req)
		assert.NoError(t, err)
		assert.NotEmpty(t, date)
	})
}
