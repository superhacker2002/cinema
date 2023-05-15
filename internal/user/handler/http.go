package handler

import (
	userRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/user/repository"
	"github.com/gorilla/mux"
	"net/http"
)

type auth interface {
	Authenticate(username string, passwordHash string) (token string, err error)
	VerifyToken(token string) (userID int, err error)
}

type httpHandler struct {
	r userRepository.Repository
}

func New(router *mux.Router, repository userRepository.Repository) httpHandler {
	handler := httpHandler{r: repository}
	handler.setRoutes(router)

	return handler
}

func (h httpHandler) setRoutes(router *mux.Router) {
	router.HandleFunc("/auth/", h.loginHandler).Methods("POST")
	s := router.PathPrefix("/users").Subrouter()
	s.HandleFunc("/", h.getUsersHandler).Methods("GET")
	s.HandleFunc("/{userId}/", h.getUserHandler).Methods("GET")
	s.HandleFunc("/", h.createUserHandler).Methods("POST")
	s.HandleFunc("/{userId}/", h.deleteUserHandler).Methods("DELETE")
	s.HandleFunc("/{userId}/", h.updateUserHandler).Methods("PUT")
}

func (h httpHandler) loginHandler(w http.ResponseWriter, _ *http.Request) {
	// TODO: login user using JWT authentication service
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
