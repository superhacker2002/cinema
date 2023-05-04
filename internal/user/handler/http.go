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
	router.HandleFunc("/auth/login", h.loginHandler).Methods("POST")
	router.HandleFunc("/users", h.getUsersHandler).Methods("GET")
	router.HandleFunc("/users", h.createUsersHandler).Methods("POST")
	router.HandleFunc("/users", h.deleteUsersHandler).Methods("DELETE")
	router.HandleFunc("/users", h.updateUsersHandler).Methods("PUT")
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
	userId := r.URL.Query()
	if _, ok := userId["userId"]; ok {
		// return user by id (only for admins)
	} else {
		// return all users (only for admins)
	}
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) createUsersHandler(w http.ResponseWriter, r *http.Request) {
	// create new user
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) deleteUsersHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query()
	if _, ok := userId["userId"]; !ok {
		http.Error(w, "user id not provided: ", http.StatusNotFound)
		return
	}
	// delete user (only for admins)
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) updateUsersHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query()
	if _, ok := userId["userId"]; !ok {
		http.Error(w, "user id not provided: ", http.StatusNotFound)
		return
	}
	// update user information
	w.WriteHeader(http.StatusOK)
}
