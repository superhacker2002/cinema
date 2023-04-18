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
	http.HandleFunc("/auth", h.handlerFunc)
	http.HandleFunc("/clients", h.handlerFunc)
	http.HandleFunc("/films", h.handlerFunc)
	http.HandleFunc("/halls", h.handlerFunc)

}

func (h httpHandler) handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
