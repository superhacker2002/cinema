package main

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/auth"
	cinemaHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinema/handler"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/config"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/repository"
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

	db, err := sql.Open("postgres", config.Db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repository := repository.New(db)

	router := mux.NewRouter()
	authentication := auth.New(config.JWTSecret)
	userHandler.New(router, authentication)
	cinemaHandler.New(router)

	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
