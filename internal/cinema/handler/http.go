package handler

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
	http.HandleFunc("/movies", moviesHandler)
	http.HandleFunc("/movies/watched", watchedMoviesHandler)
	http.HandleFunc("/halls", hallsHandler)
	http.HandleFunc("/cinema-sessions", sessionsHandler)
	http.HandleFunc("/tickets", ticketsHandler)
}

func moviesHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query()
	if _, ok := userId["movieId"]; ok {
		if r.Method == "GET" {
			// return movie by id
		} else if r.Method == "PUT" {
			// update movie information (only for admins)
		} else if r.Method == "DELETE" {
			// delete movie (only for admins)
		}
	} else {
		if r.Method == "POST" {
			// create new movie (only for admins)
		} else if r.Method == "GET" {
			// return all movies
		}
	}
	w.WriteHeader(http.StatusOK)
}

func watchedMoviesHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query()
	if _, ok := userId["userId"]; !ok {
		http.Error(w, "user id not provided: ", http.StatusBadRequest)
		return
	}

	// return list of watched movies by user id
}

func hallsHandler(w http.ResponseWriter, r *http.Request) {
	hallId := r.URL.Query()
	if _, ok := hallId["hallId"]; ok {
		if r.Method == "GET" {
			// return hall by id
		} else if r.Method == "PUT" {
			// update hall information (only for admins)
		} else if r.Method == "DELETE" {
			// delete hall (only for admins)
		}
	} else {
		if r.Method == "POST" {
			// create new hall (only for admins)
		} else if r.Method == "GET" {
			// return all halls
		}
	}
	w.WriteHeader(http.StatusOK)
}

func sessionsHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query()
	if _, ok := sessionId["sessionId"]; ok {
		if r.Method == "GET" {
			// return session by id
		} else if r.Method == "PUT" {
			// update session information (only for admins)
		} else if r.Method == "DELETE" {
			// delete session (only for admins)
		}
	} else {
		if r.Method == "POST" {
			// create new session (only for admins)
		} else if r.Method == "GET" {
			// return all sessions
		}
	}
	w.WriteHeader(http.StatusOK)
}

func ticketsHandler(w http.ResponseWriter, r *http.Request) {
	ticketId := r.URL.Query()
	if _, ok := ticketId["ticketId"]; ok {
		if r.Method == "GET" {
			// return ticket by id
		}
	} else if _, ok := ticketId["userId"]; ok {
		if r.Method == "GET" {
			// return all tickets purchased by the user
		}
	} else {
		if r.Method == "POST" {
			// create("buy") a new ticket
		}
	}
	w.WriteHeader(http.StatusOK)
}
