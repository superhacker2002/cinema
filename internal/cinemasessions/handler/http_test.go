package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/entity"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/service"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type mockService struct {
	sessions  []entity.CinemaSession
	hallId    int
	sessionId int
	err       error
}

func (m *mockService) AllSessions(date string, offset, limit int) ([]entity.CinemaSession, error) {
	return m.sessions, m.err
}

func (m *mockService) SessionsForHall(hallId int, date string) ([]entity.CinemaSession, error) {
	if hallId != m.hallId {
		return nil, service.ErrCinemaSessionsNotFound
	}
	return m.sessions, m.err
}

func (m *mockService) DeleteSession(id int) error {
	return m.err
}

func (m *mockService) CreateSession(movieId, hallId int, startTime string, price float32) (int, error) {
	return m.sessionId, m.err
}

func TestGetSessionsHandler(t *testing.T) {
	s := mockService{}
	t.Run("successful sessions get", func(t *testing.T) {
		start, _ := time.Parse(timestampLayout, "2024-05-18 20:00:00 +04")
		end, _ := time.Parse(timestampLayout, "2024-05-18 22:00:00 +04")
		session := []entity.CinemaSession{
			{
				Id:        1,
				MovieId:   1,
				StartTime: start,
				EndTime:   end,
				Status:    "scheduled",
			},
		}
		s.sessions = session
		s.hallId = 1
		s.err = nil

		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/1", nil)
		req = mux.SetURLVars(req, map[string]string{"hallId": "1"})
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: &s}.getSessionsHandler
		handler(response, req)

		assert.NotEmpty(t, response.Body.String())
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("invalid hall id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/0", nil)
		req = mux.SetURLVars(req, map[string]string{"hallId": "0"})
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: &s}.getSessionsHandler
		handler(response, req)

		assert.Equal(t, fmt.Sprintf("%v\n", ErrInvalidHallId), response.Body.String())
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("service error", func(t *testing.T) {
		s.err = service.ErrInternalError
		s.hallId = 2

		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/2", nil)
		req = mux.SetURLVars(req, map[string]string{"hallId": "2"})
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: &s}.getSessionsHandler
		handler(response, req)

		assert.Equal(t, fmt.Sprintf("%v\n", service.ErrInternalError), response.Body.String())
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	t.Run("invalid date", func(t *testing.T) {
		s.err = service.ErrInternalError
		s.hallId = 2

		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/2?date=123", nil)
		req = mux.SetURLVars(req, map[string]string{"hallId": "2"})
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: &s}.getSessionsHandler
		handler(response, req)

		assert.Equal(t, fmt.Sprintf("%v: 123\n", ErrInvalidDate), response.Body.String())
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestGetAllSessionsHandler(t *testing.T) {
	s := mockService{}
	t.Run("successful sessions get", func(t *testing.T) {
		start, _ := time.Parse(timestampLayout, "2024-05-18 20:00:00 +04")
		end, _ := time.Parse(timestampLayout, "2024-05-18 22:00:00 +04")
		sessions := []entity.CinemaSession{
			{
				Id:        1,
				MovieId:   1,
				StartTime: start,
				EndTime:   end,
				Status:    "scheduled",
			},
		}
		s.sessions = sessions
		s.err = nil

		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/", nil)
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: &s}.getAllSessionsHandler
		handler(response, req)

		assert.NotEmpty(t, response.Body.String())
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("no cinema sessions", func(t *testing.T) {
		sessions := []entity.CinemaSession{{}}
		s.sessions = sessions
		s.err = service.ErrCinemaSessionsNotFound

		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/", nil)
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: &s}.getAllSessionsHandler
		handler(response, req)

		assert.Equal(t, fmt.Sprintf("%v for all halls\n", service.ErrCinemaSessionsNotFound), response.Body.String())
		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("service error", func(t *testing.T) {
		sessions := []entity.CinemaSession{{}}
		s.sessions = sessions
		s.err = errors.New("something went wrong")

		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/", nil)
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: &s}.getAllSessionsHandler
		handler(response, req)

		assert.Equal(t, fmt.Sprintf("%v\n", service.ErrInternalError), response.Body.String())
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestPage(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "resource?offset=10&limit=20", nil)
		require.NoError(t, err, "failed to create test request")
		page, err := page(req)
		assert.NoError(t, err)
		assert.Equal(t, 10, page.offset)
		assert.Equal(t, 20, page.limit)
	})

	t.Run("invalid request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "resource?offset=-10&limit=20", nil)
		require.NoError(t, err, "failed to create test request")
		page, err := page(req)
		assert.Error(t, err)
		assert.Equal(t, 20, page.limit)
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
		assert.Equal(t, 0, page.offset)
		assert.Equal(t, 10, page.limit)
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
		assert.Equal(t, date, "2023/01/05")
	})

	t.Run("missing date", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "resource/", nil)
		require.NoError(t, err, "failed to create test request")
		date, err := date(req)
		assert.NoError(t, err)
		assert.NotEmpty(t, date)
	})
}

type errorReader struct{}

func (e errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func TestCreateSessionHandler(t *testing.T) {
	s := mockService{}
	t.Run("successful session creation", func(t *testing.T) {
		s.sessionId = 1
		s.err = nil
		hallId := "1"

		request := fmt.Sprintf(`{"movieId": %d, "startTime": "%s", "price": %f}`, 1,
			"2024-05-18 20:00:00 +04", 10.5)
		body := strings.NewReader(request)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("cinema-sessions/%s", hallId), body)
		req = mux.SetURLVars(req, map[string]string{"hallId": hallId})
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: &s}.createSessionHandler
		handler(response, req)

		assert.Equal(t, http.StatusOK, response.Code)

		var responseBody map[string]int
		err = json.Unmarshal(response.Body.Bytes(), &responseBody)
		require.NoError(t, err, "failed to parse response body")

		assert.Contains(t, responseBody, "session_id")
		assert.NotZero(t, responseBody["session_id"])
	})

	t.Run("invalid hall id", func(t *testing.T) {
		hallId := "0"

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("cinema-sessions/%s", hallId), nil)
		req = mux.SetURLVars(req, map[string]string{"hallId": hallId})
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: &s}.createSessionHandler
		handler(response, req)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, fmt.Sprintf("%s\n", ErrInvalidHallId), response.Body.String())
	})

	t.Run("failed to read request body", func(t *testing.T) {
		hallId := "1"

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("cinema-sessions/%s", hallId), errorReader{})
		req = mux.SetURLVars(req, map[string]string{"hallId": hallId})
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: &s}.createSessionHandler
		handler(response, req)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, ErrReadRequestFail.Error()+"\n", response.Body.String())
	})

	t.Run("hall is busy", func(t *testing.T) {
		hallId := "1"
		s.err = service.ErrHallIsBusy

		request := fmt.Sprintf(`{"movieId": %d, "startTime": "%s", "price": %f}`,
			1, "2024-05-18 20:00:00+4", 10.5)
		body := strings.NewReader(request)
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("cinema-sessions/%s", hallId), body)
		req = mux.SetURLVars(req, map[string]string{"hallId": hallId})
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: &s}.createSessionHandler
		handler(response, req)

		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, fmt.Sprintf("%s\n", service.ErrHallIsBusy), response.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		hallId := "1"

		request := fmt.Sprintf(`{"movieId": %d, "startTime": "%s", "price": %f}`,
			1, "2024-05-18 20:00:00+4", 10.5)
		body := strings.NewReader(request)

		s.err = errors.New("something went wrong")

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/cinema-sessions/%s", hallId), body)
		req = mux.SetURLVars(req, map[string]string{"hallId": hallId})
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()

		handler := HttpHandler{s: &s}.createSessionHandler
		handler(response, req)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, service.ErrInternalError.Error()+"\n", response.Body.String())
	})
}

//func TestDeleteSessionHandler(t *testing.T) {
//	s := mockRepo{}
//
//	t.Run("successful session deletion", func(t *testing.T) {
//		sessionID := 1
//		s.sessionId = sessionID
//		s.err = nil
//		s.sessionExists = true
//		s.serviceErr = nil
//
//		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/cinema-sessions/%d", sessionID), nil)
//		req = mux.SetURLVars(req, map[string]string{"sessionId": strconv.Itoa(sessionID)})
//		require.NoError(t, err, "failed to create test request")
//
//		response := httptest.NewRecorder()
//		s := service.New(&s)
//		handler := HttpHandler{s: s}.deleteSessionHandler
//		handler(response, req)
//
//		assert.Equal(t, "session was deleted successfully", response.Body.String())
//		assert.Equal(t, http.StatusOK, response.Code)
//	})
//
//	t.Run("invalid session ID", func(t *testing.T) {
//		sessionID := "invalid"
//		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/cinema-sessions/%s", sessionID), nil)
//		req = mux.SetURLVars(req, map[string]string{"sessionId": sessionID})
//		require.NoError(t, err, "failed to create test request")
//
//		response := httptest.NewRecorder()
//		s := service.New(&s)
//		handler := HttpHandler{s: s}.deleteSessionHandler
//		handler(response, req)
//
//		assert.Equal(t, fmt.Sprintf("%v\n", ErrInvalidSessionId), response.Body.String())
//		assert.Equal(t, http.StatusBadRequest, response.Code)
//	})
//
//	t.Run("session not found", func(t *testing.T) {
//		sessionID := 2
//		s.sessionId = sessionID
//		s.sessionExists = false
//
//		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/cinema-sessions/%d", sessionID), nil)
//		req = mux.SetURLVars(req, map[string]string{"sessionId": strconv.Itoa(sessionID)})
//		require.NoError(t, err, "failed to create test request")
//
//		response := httptest.NewRecorder()
//		s := service.New(&s)
//		handler := HttpHandler{s: s}.deleteSessionHandler
//		handler(response, req)
//
//		assert.Equal(t, fmt.Sprintf("%v\n", service.ErrCinemaSessionsNotFound), response.Body.String())
//		assert.Equal(t, http.StatusNotFound, response.Code)
//	})
//
//	t.Run("ssitory error", func(t *testing.T) {
//		sessionID := 3
//		s.sessionId = sessionID
//		s.err = errors.New("something went wrong")
//		s.sessionExists = true
//		s.serviceErr = nil
//
//		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/cinema-sessions/%d", sessionID), nil)
//		req = mux.SetURLVars(req, map[string]string{"sessionId": strconv.Itoa(sessionID)})
//		require.NoError(t, err, "failed to create test request")
//
//		response := httptest.NewRecorder()
//		s := service.New(&s)
//		handler := HttpHandler{s: s}.deleteSessionHandler
//		handler(response, req)
//
//		assert.Equal(t, fmt.Sprintf("%v\n", service.ErrInternalError), response.Body.String())
//		assert.Equal(t, http.StatusInternalServerError, response.Code)
//	})
//}
