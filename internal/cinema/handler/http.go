package handler

import (
	"net/http"
)

type httpHandler struct{}

func New() httpHandler {
	handler := httpHandler{}
	handler.setRoutes()

	return handler
}

// setRoutes adds handlers to http.DefaultServeMux
func (h httpHandler) setRoutes() {
	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	http.HandleFunc("/movies", handler)
	http.HandleFunc("/halls", handler)
}
