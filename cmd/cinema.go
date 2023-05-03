package main

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/auth"
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
	authentication := auth.New(config.JWTSecret)
	userHandler.New(router, authentication)
	cinemaHandler.New(router)

	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
