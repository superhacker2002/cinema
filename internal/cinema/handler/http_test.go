package handler

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetRoutes(t *testing.T) {
	handler := HttpHandler{}
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
		{path: "/halls/", status: http.StatusOK, method: "POST"},
		{path: "/halls/", status: http.StatusOK, method: "GET"},
		{path: "/halls/1/", status: http.StatusOK, method: "GET"},
		{path: "/halls/1/", status: http.StatusOK, method: "PUT"},
		{path: "/halls/1/", status: http.StatusOK, method: "DELETE"},

		{path: "/movies/", status: http.StatusOK, method: "POST"},
		{path: "/movies/", status: http.StatusOK, method: "GET"},
		{path: "/movies/1/", status: http.StatusOK, method: "GET"},
		{path: "/movies/1/", status: http.StatusOK, method: "PUT"},
		{path: "/movies/1/", status: http.StatusOK, method: "DELETE"},
		{path: "/movies/watched/1/", status: http.StatusOK, method: "GET"},

		{path: "/cinema-sessions/", status: http.StatusOK, method: "POST"},
		{path: "/cinema-sessions/", status: http.StatusOK, method: "GET"},
		{path: "/cinema-sessions/1/", status: http.StatusOK, method: "GET"},
		{path: "/cinema-sessions/1/", status: http.StatusOK, method: "PUT"},
		{path: "/cinema-sessions/1/", status: http.StatusOK, method: "DELETE"},

		{path: "/tickets/1/", status: http.StatusOK, method: "GET"},
		{path: "/tickets/", status: http.StatusOK, method: "POST"},

		{path: "/invalid/", status: http.StatusNotFound, method: "GET"},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest(tc.method, server.URL+tc.path, nil)
		assert.NoError(t, err)

		resp, clientErr := client.Do(req)
		assert.NoError(t, clientErr)
		assert.Equal(t, tc.status, resp.StatusCode, "Request to %s, method %s", tc.path, tc.method)
	}
}
