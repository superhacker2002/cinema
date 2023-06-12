package handler

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/apiutils"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/hall/service"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	ErrReadRequestFail = errors.New("failed to read request")
	ErrInvalidHallId   = errors.New("invalid hall id")
)

type cinemaHall struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

type Service interface {
	Halls() ([]service.Hall, error)
	HallById(id int) (service.Hall, error)
	CreateHall(name string, capacity int) (hallId int, err error)
	UpdateHall(id int, name string, capacity int) (err error)
	DeleteHall(id int) error
}

type AccessChecker interface {
	Authenticate(next http.Handler) http.Handler
	CheckPerms(perms ...string) mux.MiddlewareFunc
}

type HttpHandler struct {
	s Service
}

func New(s Service) HttpHandler {
	return HttpHandler{
		s: s,
	}
}

func (h HttpHandler) SetRoutes(router *mux.Router, a AccessChecker) {
	userRouter := router.PathPrefix("/halls").Subrouter()
	userRouter.Use(a.Authenticate)

	userRouter.HandleFunc("/", h.getHallsHandler).Methods(http.MethodGet)
	userRouter.HandleFunc("/{hallId}", h.getHallHandler).Methods(http.MethodGet)

	adminRouter := router.PathPrefix("/halls").Subrouter()
	adminRouter.Use(a.Authenticate)
	adminRouter.Use(a.CheckPerms(service.AdminRole))

	adminRouter.HandleFunc("/", h.createHallHandler).Methods(http.MethodPost)
	adminRouter.HandleFunc("/{hallId}", h.updateHallHandler).Methods(http.MethodPut)
	adminRouter.HandleFunc("/{hallId}", h.deleteHallHandler).Methods(http.MethodDelete)
}

func (h HttpHandler) getHallsHandler(w http.ResponseWriter, _ *http.Request) {
	cinemaHalls, err := h.s.Halls()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, entitiesToDTO(cinemaHalls), http.StatusOK)
}

func (h HttpHandler) createHallHandler(w http.ResponseWriter, r *http.Request) {
	type hallInfo struct {
		Name     string `json:"name"`
		Capacity int    `json:"capacity"`
	}

	var hall hallInfo

	err := json.NewDecoder(r.Body).Decode(&hall)
	if err != nil {
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.s.CreateHall(hall.Name, hall.Capacity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, map[string]int{"hallId": id}, http.StatusCreated)
}

func (h HttpHandler) getHallHandler(w http.ResponseWriter, r *http.Request) {
	hallID, err := apiutils.IntPathParam(r, "hallId")
	if err != nil {
		http.Error(w, ErrInvalidHallId.Error(), http.StatusBadRequest)
		return
	}

	hall, err := h.s.HallById(hallID)
	if errors.Is(err, service.ErrHallNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiutils.WriteResponse(w, entityToDTO(hall), http.StatusOK)
}

func (h HttpHandler) updateHallHandler(w http.ResponseWriter, r *http.Request) {
	hallID, err := apiutils.IntPathParam(r, "hallId")
	if err != nil {
		http.Error(w, ErrInvalidHallId.Error(), http.StatusBadRequest)
		return
	}

	type hallInfo struct {
		Name     string `json:"name"`
		Capacity int    `json:"capacity"`
	}

	var hall hallInfo

	err = json.NewDecoder(r.Body).Decode(&hall)
	if err != nil {
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.UpdateHall(hallID, hall.Name, hall.Capacity)
	if errors.Is(err, service.ErrHallNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h HttpHandler) deleteHallHandler(w http.ResponseWriter, r *http.Request) {
	hallID, err := apiutils.IntPathParam(r, "hallId")
	if err != nil {
		http.Error(w, ErrInvalidHallId.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.DeleteHall(hallID)
	if errors.Is(err, service.ErrHallNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func entitiesToDTO(halls []service.Hall) []cinemaHall {
	var DTOHalls []cinemaHall
	for _, hall := range halls {
		DTOHalls = append(DTOHalls, entityToDTO(hall))
	}
	return DTOHalls
}

func entityToDTO(hall service.Hall) cinemaHall {
	return cinemaHall{
		ID:       hall.Id,
		Name:     hall.Name,
		Capacity: hall.Capacity,
	}
}
