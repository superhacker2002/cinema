package user

import (
	"net/http"
)

type httpHandler struct{}

func New() httpHandler {
	handler := httpHandler{}
	handler.setRoutes()

	return handler
}

// setRoutes adds handlers to http.DefaultServeMux
func (h httpHandler) setRoutes() {
	http.HandleFunc("/auth/login", logIn)
	http.HandleFunc("/users", userHandler)
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
