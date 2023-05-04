package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

type httpHandler struct{}

func New(router *mux.Router) httpHandler {
	handler := httpHandler{}
	handler.setRoutes(router)

	return handler
}

func (h httpHandler) setRoutes(router *mux.Router) {
	router.HandleFunc("/auth/login", h.logIn)
	router.HandleFunc("/users", h.getUsersHandler).Methods("GET")
	router.HandleFunc("/users/{userId}", h.getUserHandler).Methods("GET")
	router.HandleFunc("/users", h.createUserHandler).Methods("POST")
	router.HandleFunc("/users/{userId}", h.deleteUserHandler).Methods("DELETE")
	router.HandleFunc("/users/{userId}", h.updateUserHandler).Methods("PUT")

}

func (h httpHandler) logIn(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return all users
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) getUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return user by ID
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) createUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: create user
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: delete user by ID
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: update user by ID
	w.WriteHeader(http.StatusOK)
}
