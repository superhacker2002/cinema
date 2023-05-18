package main

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/auth"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/config"

	userHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/user/handler"
	userRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/user/repository"

	sessionsHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/handler"
	sessionsRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/repository"

	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Lshortfile)
	configs, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	if err = configs.Validate(); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", configs.Db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	userRepo := userRepository.New(db)
	sessionsRepo := sessionsRepository.New(db)

	router := mux.NewRouter()
	authentication := auth.New(configs.JWTSecret, configs.TokenExp, userRepo)

	userHandler.New(router, authentication, userRepo)
	sessionsHandler.New(router, sessionsRepo)

	log.Fatal(http.ListenAndServe(":"+configs.Port, router))
}
