package main

import (
	cinemaHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinema/handler"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/config"
	userHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/user/handler"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	config := config.New()

	router := mux.NewRouter()
	userHandler.New(router)
	cinemaHandler.New(router)

	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
