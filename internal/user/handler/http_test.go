package handler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetRoutes(t *testing.T) {
	handler := httpHandler{}
	handler.setRoutes()

	server := httptest.NewServer(http.DefaultServeMux)
	defer server.Close()

	client := server.Client()

	testCases := []struct {
		path   string
		status int
	}{
		{path: "/auth/login", status: http.StatusOK},
		{path: "/users", status: http.StatusOK},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest("GET", server.URL+tc.path, nil)
		assert.NoError(t, err)

		resp, clientErr := client.Do(req)
		assert.NoError(t, clientErr)
		assert.Equal(t, tc.status, resp.StatusCode, "Request to %s", tc.path)
	}
}
