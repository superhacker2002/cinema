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
	ErrReadRequestFail = errors.New("failed to read request body")
	ErrNoUsername      = errors.New("missing username")
	ErrNoPassword      = errors.New("missing password")
	ErrInternalError   = errors.New("internal server error")
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type auth interface {
	Authenticate(username string, passwordHash string) (token string, err error)
	VerifyToken(token string) (userID int, err error)
}

type HttpHandler struct {
	a auth
	r userRepository.Repository
}

func New(router *mux.Router, auth auth, repo userRepository.Repository) HttpHandler {
	handler := HttpHandler{a: auth, r: repo}
	handler.setRoutes(router)

	return handler
}

func (h HttpHandler) setRoutes(router *mux.Router) {
	router.HandleFunc("/auth/", h.loginHandler).Methods("POST")
	s := router.PathPrefix("/users").Subrouter()
	s.HandleFunc("/", h.createUserHandler).Methods("POST")
	s.HandleFunc("/", h.getUsersHandler).Methods("GET")
	s.HandleFunc("/{userId}/", h.getUserHandler).Methods("GET")
	s.HandleFunc("/{userId}/", h.deleteUserHandler).Methods("DELETE")
	s.HandleFunc("/{userId}/", h.updateUserHandler).Methods("PUT")
}

func (h HttpHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
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

	if err = creds.validate(); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	if c.Username == "" {
		return ErrNoUsername
	} else if c.Password == "" {
		return ErrNoPassword
	}
	return nil
}

func (h HttpHandler) createUserHandler(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.r.CreateUser(creds.Username, creds.Password)
	if errors.Is(err, userRepository.ErrUserExists) {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"user_id": id})
}

func (h HttpHandler) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return all users
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) getUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: return user by ID
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: delete user by ID
	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: update user by ID
	w.WriteHeader(http.StatusOK)
}
