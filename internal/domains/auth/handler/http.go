package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/apiutils"
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

type service interface {
	Authenticate(username string, passwordHash string) (token string, err error)
}

type HttpHandler struct {
	s service
}

func New(router *mux.Router, s service) HttpHandler {
	handler := HttpHandler{s: s}
	handler.setRoutes(router)

	return handler
}

func (h HttpHandler) setRoutes(router *mux.Router) {
	router.HandleFunc("/auth/", h.loginHandler).Methods("POST")
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

	token, err := h.s.Authenticate(creds.Username, creds.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to authenticate: "+err.Error(), http.StatusUnauthorized)
		return
	}

	apiutils.WriteResponse(w, map[string]string{"token": token}, http.StatusCreated)
}

func (c credentials) validate() error {
	if c.Username == "" {
		return ErrNoUsername
	} else if c.Password == "" {
		return ErrNoPassword
	}
	return nil
}
