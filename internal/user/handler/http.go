package handler

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

var (
	ErrReadRequestFail      = errors.New("failed to read request body")
	ErrNoUsernameOrPassword = errors.New("missing username or password")
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type repository interface {
	CreateUser(username string, password string) (string, error)
}

type auth interface {
	Authenticate(username string, password string) (string, error)
	VerifyToken(tokenString string) (string, error)
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
	s.HandleFunc("/", h.createUserHandler).Methods("POST")
	s.HandleFunc("/", h.getUsersHandler).Methods("GET")
	s.HandleFunc("/{userId}/", h.getUserHandler).Methods("GET")
	s.HandleFunc("/{userId}/", h.deleteUserHandler).Methods("DELETE")
	s.HandleFunc("/{userId}/", h.updateUserHandler).Methods("PUT")
}

func (h httpHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	var creds credentials
	err = json.Unmarshal(body, &creds)
	if err != nil {
		log.Println(err)
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	if creds.Username == "" || creds.Password == "" {
		log.Println(err)
		http.Error(w, ErrNoUsernameOrPassword.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.auth.Authenticate(creds.Username, creds.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to authenticate: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h httpHandler) createUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	var creds credentials
	err = json.Unmarshal(body, &creds)
	if err != nil {
		log.Println(err)
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	if creds.Username == "" || creds.Password == "" {
		log.Println(err)
		http.Error(w, ErrNoUsernameOrPassword.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.repository.CreateUser(creds.Username, creds.Password)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"user_id": id})
}

func (h httpHandler) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return all users
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) getUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return user by ID
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
