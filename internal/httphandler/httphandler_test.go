package httphandler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetRoutes(t *testing.T) {
	handler := httpHandler{}
	handler.setRoutes()

	server := httptest.Server{}
	defer server.Close()

	for _, route := range []string{"/auth", "/clients", "/films", "/halls"} {
		resp, err := http.Get(server.URL + route)
		if err != nil {
			t.Errorf("Error making HTTP request to %s: %s", route, err.Error())
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected HTTP status code 200, but got %d", resp.StatusCode)
		}
	}
}
