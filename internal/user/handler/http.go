package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

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
	router.HandleFunc("/auth/login/", h.loginHandler)
	s := router.PathPrefix("/users").Subrouter()
	s.HandleFunc("/", h.getUsersHandler).Methods("GET")
	s.HandleFunc("/{userId}/", h.getUserHandler).Methods("GET")
	s.HandleFunc("/", h.createUserHandler).Methods("POST")
	s.HandleFunc("/{userId}/", h.deleteUserHandler).Methods("DELETE")
	s.HandleFunc("/{userId}/", h.updateUserHandler).Methods("PUT")
}

func (h httpHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	var creds credentials
	log.Println(string(body))
	err = json.Unmarshal(body, &creds)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(creds.Username, creds.Password)

	token, err := h.auth.Authenticate(creds.Username, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
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
