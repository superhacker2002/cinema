package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinema/repository"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockRepo struct {
	session repository.CinemaSession
	hallId  int
	err     error
}

func (m *mockRepo) SessionsForHall(hallId int, currentTime string) ([]repository.CinemaSession, error) {
	var cinemaSessions []repository.CinemaSession
	if hallId != m.hallId {
		return nil, repository.ErrCinemaSessionsNotFound
	}
	cinemaSessions = append(cinemaSessions, m.session)
	return cinemaSessions, m.err
}

func TestGetSessionsHandler(t *testing.T) {
	repo := mockRepo{}
	t.Run("successful sessions get", func(t *testing.T) {
		session := repository.CinemaSession{
			ID:        1,
			MovieId:   1,
			StartTime: "2023-05-18 20:00:00",
			EndTime:   "2023-05-18 22:00:00",
			Status:    "scheduled",
		}
		repo.session = session
		repo.hallId = 1
		repo.err = nil

		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/1", nil)
		req = mux.SetURLVars(req, map[string]string{"hallId": "1"})
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{r: &repo}.getSessionsHandler
		handler(response, req)

		assert.Equal(t, "[{\"ID\":1,\"MovieId\":1,\"StartTime\":\"2023-05-18 20:00:00\","+
			"\"EndTime\":\"2023-05-18 22:00:00\",\"Status\":\"scheduled\"}]\n", response.Body.String())
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("invalid hall id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/0", nil)
		req = mux.SetURLVars(req, map[string]string{"hallId": "0"})
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{r: &repo}.getSessionsHandler
		handler(response, req)

		assert.Equal(t, ErrInvalidHallId.Error()+": 0\n", response.Body.String())
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("no cinema sessions", func(t *testing.T) {
		session := repository.CinemaSession{}
		repo.session = session
		repo.hallId = 1
		repo.err = repository.ErrCinemaSessionsNotFound

		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/2", nil)
		req = mux.SetURLVars(req, map[string]string{"hallId": "2"})
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{r: &repo}.getSessionsHandler
		handler(response, req)

		assert.Equal(t, repository.ErrCinemaSessionsNotFound.Error()+" for hall 2\n", response.Body.String())
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		session := repository.CinemaSession{}
		repo.session = session
		repo.hallId = 2
		repo.err = errors.New("something went wrong")

		req, err := http.NewRequest(http.MethodGet, "cinema-sessions/2", nil)
		req = mux.SetURLVars(req, map[string]string{"hallId": "2"})
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{r: &repo}.getSessionsHandler
		handler(response, req)

		assert.Equal(t, ErrInternalError.Error()+"\n", response.Body.String())
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}
