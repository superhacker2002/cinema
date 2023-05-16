package handler

import (
	userRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/user/repository"
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

type auth interface {
	Authenticate(username string, passwordHash string) (token string, err error)
	VerifyToken(token string) (userID int, err error)
}

type httpHandler struct {
	a auth
	r userRepository.Repository
}

func New(router *mux.Router, auth auth, repo userRepository.Repository) httpHandler {
	handler := httpHandler{a: auth, r: repo}
	handler.setRoutes(router)

	return handler
}

func (h httpHandler) setRoutes(router *mux.Router) {
	router.HandleFunc("/auth/", h.loginHandler).Methods("POST")
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

	token, err := h.a.Authenticate(creds.Username, creds.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to authenticate: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (c credentials) validate() error {
	if c.Username == "" || c.Password == "" {
		return ErrNoUsernameOrPassword
	}
	return nil
}

func (h httpHandler) createUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	var creds credentials

	if err = json.Unmarshal(body, &creds); err != nil {
		log.Println(err)
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	if err = creds.validate(); err != nil {
		log.Println(err)
		http.Error(w, ErrNoUsernameOrPassword.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.r.CreateUser(creds.Username, creds.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to sign up: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"user_id": id})
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
