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
	router.HandleFunc("/auth/login", logIn)
	router.HandleFunc("/users", userHandler)
}

func logIn(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query()
	if _, ok := userId["userId"]; ok {
		if r.Method == "GET" {
			// return user by id (only for admins)
		} else if r.Method == "PUT" {
			// update user information
		} else if r.Method == "DELETE" {
			// delete user (only for admins)
		}
	} else {
		if r.Method == "POST" {
			// create new user
		} else if r.Method == "GET" {
			// return all users (only for admins)
		}
	}
	w.WriteHeader(http.StatusOK)
}
