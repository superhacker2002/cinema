package handler

import (
	userRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/user/repository"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSetRoutes(t *testing.T) {
	handler := httpHandler{}
	router := mux.NewRouter()
	handler.setRoutes(router)

	server := httptest.NewServer(router)
	defer server.Close()

	client := server.Client()

	testCases := []struct {
		path   string
		status int
		method string
	}{
		{path: "/users/", status: http.StatusOK, method: "GET"},
		{path: "/users/1/", status: http.StatusOK, method: "GET"},
		{path: "/users/1/", status: http.StatusOK, method: "PUT"},
		{path: "/users/1/", status: http.StatusOK, method: "DELETE"},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest(tc.method, server.URL+tc.path, nil)
		assert.NoError(t, err)

		resp, clientErr := client.Do(req)
		assert.NoError(t, clientErr)
		assert.Equal(t, tc.status, resp.StatusCode, "Request to %s using method %s", tc.path, tc.method)
	}
}

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
		handler := httpHandler{a: auth}.loginHandler
		handler(response, req)

		assert.Equal(t, "{\"token\":\"test_token\"}\n", response.Body.String())
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("reading request fail", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "auth/",
			strings.NewReader(`invalid json`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := httpHandler{a: auth}.loginHandler
		handler(response, req)

		assert.Equal(t, "failed to read request body\n", response.Body.String())
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("authentication fail", func(t *testing.T) {
		auth.err = errors.New("something went wrong")
		req, err := http.NewRequest(http.MethodPost, "auth/",
			strings.NewReader(`{"username": "test_user", "password": "test_password"}`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := httpHandler{a: auth}.loginHandler
		handler(response, req)

		assert.Equal(t, "failed to authenticate: something went wrong\n", response.Body.String())
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	})
}

type mockRepository struct {
	err    error
	userId int
}

func (m mockRepository) GetUser(username string) (userRepository.Credentials, error) {
	return userRepository.Credentials{}, nil
}

func (m mockRepository) CreateUser(username string, passwordHash string) (int, error) {
	return m.userId, m.err
}

func TestCreateUserHandler(t *testing.T) {
	auth := mockAuth{}
	repo := mockRepository{}
	t.Run("successful registration", func(t *testing.T) {
		repo.err = nil
		repo.userId = 1
		req, err := http.NewRequest(http.MethodPost, "users/",
			strings.NewReader(`{"username": "test_user", "password": "test_password"}`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := httpHandler{a: auth, r: repo}.createUserHandler
		handler(response, req)

		assert.Equal(t, "{\"user_id\":1}\n", response.Body.String())
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("reading request fail", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "users/",
			strings.NewReader(`invalid json`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := httpHandler{a: auth, r: repo}.createUserHandler
		handler(response, req)

		assert.Equal(t, "failed to read request body\n", response.Body.String())
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("missing password", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "users/",
			strings.NewReader(`{"username": "test_user"}`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := httpHandler{a: auth, r: repo}.createUserHandler
		handler(response, req)

		assert.Equal(t, "missing username or password\n", response.Body.String())
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("user creation fail", func(t *testing.T) {
		repo.err = errors.New("something went wrong")
		req, err := http.NewRequest(http.MethodPost, "users/",
			strings.NewReader(`{"username": "test_user", "password": "test_password"}`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := httpHandler{a: auth, r: repo}.createUserHandler
		handler(response, req)

		assert.Equal(t, "failed to sign up: internal server error\n", response.Body.String())
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}
