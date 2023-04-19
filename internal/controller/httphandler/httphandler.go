package httphandler

import (
	"net/http"
)

type httpHandler struct{}

func New() httpHandler {
	handler := httpHandler{}
	handler.setRoutes()
	return handler
}

func (h httpHandler) setRoutes() {
	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	http.HandleFunc("/auth", handler)
	http.HandleFunc("/clients", handler)
	http.HandleFunc("/films", handler)
	http.HandleFunc("/halls", handler)
}
