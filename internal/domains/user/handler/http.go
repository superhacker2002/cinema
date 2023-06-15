package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/apiutils"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/user/service"
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
	ErrInvalidUserId   = errors.New("invalid user id")
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type accessChecker interface {
	Authenticate(next http.Handler) http.Handler
	CheckPerms(perms ...string) mux.MiddlewareFunc
}

type Service interface {
	CreateUser(username string, passwordHash string) (userId int, err error)
	MakeAdmin(userId int) error
}

type HttpHandler struct {
	s Service
}

func New(router *mux.Router, s Service) HttpHandler {
	return HttpHandler{
		s: s,
	}
}

func (h HttpHandler) SetRoutes(router *mux.Router, a accessChecker) {
	allRouter := router.PathPrefix("/users").Subrouter()
	allRouter.HandleFunc("/", h.createUserHandler).Methods("POST")

	adminRouter := router.PathPrefix("/users").Subrouter()
	adminRouter.Use(a.Authenticate)
	adminRouter.Use(a.CheckPerms("admin"))
	allRouter.HandleFunc("/{userId}/grand-admin", h.makeUserAdmin).Methods("POST")

}

func (h HttpHandler) makeUserAdmin(w http.ResponseWriter, r *http.Request) {
	userId, err := apiutils.IntPathParam(r, "userId")
	if err != nil {
		http.Error(w, ErrInvalidUserId.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.MakeAdmin(userId)
	if errors.Is(err, service.ErrUserNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	apiutils.WriteResponse(w, map[string]int{"userId": id}, http.StatusCreated)
}

func (c credentials) validate() error {
	if c.Username == "" {
		return ErrNoUsername
	} else if c.Password == "" {
		return ErrNoPassword
	}
	return nil
}
