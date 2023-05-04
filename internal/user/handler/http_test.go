package handler

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
		{path: "/auth/login", status: http.StatusOK, method: "POST"},
		{path: "/users", status: http.StatusOK, method: "POST"},
		{path: "/users", status: http.StatusOK, method: "GET"},
		{path: "/users?userId=1", status: http.StatusOK, method: "GET"},
		{path: "/users", status: http.StatusNotFound, method: "PUT"},
		{path: "/users?userId=1", status: http.StatusOK, method: "PUT"},
		{path: "/users", status: http.StatusNotFound, method: "DELETE"},
		{path: "/users?userId=1", status: http.StatusOK, method: "DELETE"},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest(tc.method, server.URL+tc.path, nil)
		assert.NoError(t, err)

		resp, clientErr := client.Do(req)
		assert.NoError(t, clientErr)
		assert.Equal(t, tc.status, resp.StatusCode, "Request to %s using method %s", tc.path, tc.method)
	}
}
