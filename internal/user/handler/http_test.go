package handler

import (
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
		{path: "/users/", status: http.StatusOK, method: "POST"},
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
func (m mockAuth) VerifyToken(returnTokenString string) (string, error) {
	return "", nil
}

type loginResponse struct {
	Token string `json:"token"`
}

func TestLoginHandler(t *testing.T) {
	auth := mockAuth{}
	t.Run("successful authentication", func(t *testing.T) {
		auth.token = "test_token"
		auth.err = nil
		req, err := http.NewRequest(http.MethodPost, "auth/login/",
			strings.NewReader(`{"username": "test_user", "password": "test_password"}`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := httpHandler{auth: auth}.loginHandler
		handler(response, req)

		assert.NotEmpty(t, response.Body)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("reading request fail", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "auth/login/",
			strings.NewReader(`invalid json`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := httpHandler{auth: auth}.loginHandler
		handler(response, req)

		assert.Equal(t, "failed to read request body\n", response.Body.String())
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("unmarshal request fail", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "auth/login/",
			strings.NewReader(`{"username": "test_user"}`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := httpHandler{auth: auth}.loginHandler
		handler(response, req)

		assert.Equal(t, "missing username or password\n", response.Body.String())
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("authentication fail", func(t *testing.T) {
		auth.err = errors.New("something went wrong")
		req, err := http.NewRequest(http.MethodPost, "auth/login/",
			strings.NewReader(`{"username": "test_user", "password": "test_password"}`))
		require.NoError(t, err, "failed to create test request")

		response := httptest.NewRecorder()
		handler := httpHandler{auth: auth}.loginHandler
		handler(response, req)

		assert.Equal(t, "failed to authenticate: something went wrong\n", response.Body.String())
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	})
}
