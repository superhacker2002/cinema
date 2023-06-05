package handler

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockAuth struct {
	token string
	err   error
}

func (m mockAuth) Authenticate(username string, password string) (string, error) {
	return m.token, m.err
}
func (m mockAuth) VerifyToken(token string) (int, error) {
	return 0, nil
}

func TestLoginHandler(t *testing.T) {
	auth := mockAuth{}
	t.Run("successful authentication", func(t *testing.T) {
		auth.token = "test_token"
		auth.err = nil
		req, err := http.NewRequest(http.MethodPost, "auth/",
			strings.NewReader(`{"username": "test_user", "password": "test_password"}`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: auth}.loginHandler
		handler(response, req)

		assert.Equal(t, "{\"token\":\"test_token\"}\n", response.Body.String())
		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("reading request fail", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "auth/",
			strings.NewReader(`invalid json`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: auth}.loginHandler
		handler(response, req)

		assert.Equal(t, ErrReadRequestFail.Error()+"\n", response.Body.String())
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("no username provided", func(t *testing.T) {
		auth.token = "test_token"
		auth.err = nil
		req, err := http.NewRequest(http.MethodPost, "auth/",
			strings.NewReader(`{"password": "test_password"}`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: auth}.loginHandler
		handler(response, req)

		assert.Equal(t, ErrNoUsername.Error()+"\n", response.Body.String())
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("no password provided", func(t *testing.T) {
		auth.token = "test_token"
		auth.err = nil
		req, err := http.NewRequest(http.MethodPost, "auth/",
			strings.NewReader(`{"username": "test_user"}`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: auth}.loginHandler
		handler(response, req)

		assert.Equal(t, ErrNoPassword.Error()+"\n", response.Body.String())
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("authentication fail", func(t *testing.T) {
		auth.err = errors.New("something went wrong")
		req, err := http.NewRequest(http.MethodPost, "auth/",
			strings.NewReader(`{"username": "test_user", "password": "test_password"}`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := HttpHandler{s: auth}.loginHandler
		handler(response, req)

		assert.Equal(t, "failed to authenticate: something went wrong\n", response.Body.String())
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	})
}
