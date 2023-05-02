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
	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	userHandler.New(router, config.JWTSecret)
	cinemaHandler.New(router)

	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
