package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

type httpHandler struct {
	JWTSecret []byte
}

func New(router *mux.Router, JWTSecret []byte) httpHandler {
	handler := httpHandler{JWTSecret: JWTSecret}
	handler.setRoutes(router)

	return handler
}

func (h httpHandler) setRoutes(router *mux.Router) {
	router.HandleFunc("/auth/login", loginHandler).Methods("POST")
	router.HandleFunc("/users", getUsersHandler).Methods("GET")
	router.HandleFunc("/users", createUsersHandler).Methods("POST")
	router.HandleFunc("/users", deleteUsersHandler).Methods("DELETE")
	router.HandleFunc("/users", updateUsersHandler).Methods("PUT")
}

func loginHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query()
	if _, ok := userId["userId"]; ok {
		// return user by id (only for admins)
	} else {
		// return all users (only for admins)
	}
	w.WriteHeader(http.StatusOK)
}

func createUsersHandler(w http.ResponseWriter, r *http.Request) {
	// create new user
	w.WriteHeader(http.StatusOK)
}

func deleteUsersHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query()
	if _, ok := userId["userId"]; !ok {
		http.Error(w, "user id not provided: ", http.StatusNotFound)
		return
	}
	// delete user (only for admins)
	w.WriteHeader(http.StatusOK)
}

func updateUsersHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query()
	if _, ok := userId["userId"]; !ok {
		http.Error(w, "user id not provided: ", http.StatusNotFound)
		return
	}
	// update user information
	w.WriteHeader(http.StatusOK)
}
