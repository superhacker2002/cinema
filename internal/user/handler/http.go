package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/apiutils"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/user/service"
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
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Service interface {
	CreateUser(username string, passwordHash string) (userId int, err error)
}

type HttpHandler struct {
	s Service
}

func New(router *mux.Router, s Service) HttpHandler {
	handler := HttpHandler{s: s}
	handler.setRoutes(router)

	return handler
}

func (h HttpHandler) setRoutes(router *mux.Router) {
	s := router.PathPrefix("/users").Subrouter()
	s.HandleFunc("/", h.createUserHandler).Methods("POST")
	s.HandleFunc("/", h.getUsersHandler).Methods("GET")
	s.HandleFunc("/{userId}/", h.getUserHandler).Methods("GET")
	s.HandleFunc("/{userId}/", h.deleteUserHandler).Methods("DELETE")
	s.HandleFunc("/{userId}/", h.updateUserHandler).Methods("PUT")
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

	id, err := h.s.CreateUser(creds.Username, creds.Password)
	if errors.Is(err, service.ErrUserExists) {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, map[string]int{"user_id": id}, http.StatusCreated)
}

func (c credentials) validate() error {
	if c.Username == "" {
		return ErrNoUsername
	} else if c.Password == "" {
		return ErrNoPassword
	}
	return nil
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
