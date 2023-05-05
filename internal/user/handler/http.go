package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

type repository interface{}

type auth interface {
	Authenticate(username string, password string) (string, error)
}

type httpHandler struct {
	auth       auth
	repository repository
}

func New(router *mux.Router, auth auth, repository repository) httpHandler {
	handler := httpHandler{auth: auth, repository: repository}
	handler.setRoutes(router)

	return handler
}

func (h httpHandler) setRoutes(router *mux.Router) {
	router.HandleFunc("/auth/login/", h.loginHandler).Methods("POST")
	s := router.PathPrefix("/users").Subrouter()
	s.HandleFunc("/", h.getUsersHandler).Methods("GET")
	s.HandleFunc("/{userId}/", h.getUserHandler).Methods("GET")
	s.HandleFunc("/", h.createUserHandler).Methods("POST")
	s.HandleFunc("/{userId}/", h.deleteUserHandler).Methods("DELETE")
	s.HandleFunc("/{userId}/", h.updateUserHandler).Methods("PUT")
}

func (h httpHandler) loginHandler(w http.ResponseWriter, _ *http.Request) {
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
