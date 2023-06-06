package handler

import (
	pdfgenerator "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/ticket/pdfservice"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	ErrReadRequestFail = errors.New("failed to read request body")
)

type HttpHandler struct{}

type ticket struct {
	SessionId  int `json:"sessionId"`
	UserId     int `json:"userId"`
	SeatNumber int `json:"seatNumber"`
}

type AccessChecker interface {
	Authorize(next http.Handler) http.Handler
	CheckPerms(next http.Handler, perms []string) http.Handler
}

func New() HttpHandler {
	return HttpHandler{}
}

func (h HttpHandler) SetRoutes(router *mux.Router, a AccessChecker) {
	s := router.PathPrefix("/tickets").Subrouter()
	s.Use(a.Authorize)
	s.HandleFunc("/", h.createTicket).Methods(http.MethodPost)
}

func (h HttpHandler) createTicket(w http.ResponseWriter, r *http.Request) {
	err := pdfgenerator.GeneratePDF(1, 2, 3, "output.pdf")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("PDF generated successfully")
}
