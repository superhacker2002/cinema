package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/user/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockService struct {
	err    error
	userId int
}

func (m mockService) CreateUser(username string, passwordHash string) (int, error) {
	return m.userId, m.err
}

func TestCreateUserHandler(t *testing.T) {
	s := mockService{}
	t.Run("successful registration", func(t *testing.T) {
		s.err = nil
		s.userId = 1
		req, err := http.NewRequest(http.MethodPost, "users/",
			strings.NewReader(`{"username": "test_user", "password": "test_password"}`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: s}.createUserHandler
		handler(response, req)

		assert.Equal(t, "{\"user_id\":1}\n", response.Body.String())
		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("reading request fail", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "users/",
			strings.NewReader(`invalid json`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: s}.createUserHandler
		handler(response, req)

		assert.Equal(t, ErrReadRequestFail.Error()+"\n", response.Body.String())
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("missing password", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "users/",
			strings.NewReader(`{"username": "test_user"}`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: s}.createUserHandler
		handler(response, req)

		assert.Equal(t, ErrNoPassword.Error()+"\n", response.Body.String())
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("user creation fail", func(t *testing.T) {
		s.err = service.ErrInternalError
		req, err := http.NewRequest(http.MethodPost, "users/",
			strings.NewReader(`{"username": "test_user", "password": "test_password"}`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: s}.createUserHandler
		handler(response, req)

		assert.Equal(t, service.ErrInternalError.Error()+"\n", response.Body.String())
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}
