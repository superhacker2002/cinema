package main

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/auth/service"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/config"
	"os"
	"time"

	userHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/user/handler"
	userRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/user/repository"

	sessionsHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/handler"
	sessionsRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/repository"
	sessionsService "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasessions/service"

	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var timeZone = time.FixedZone("UTC+4", 4*60*60)

func main() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stdout)
	configs, err := config.New()
	if err != nil {
		log.Fatalf("config loading failed: %v", err)
	}

	db, err := sql.Open("postgres", configs.Db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	userRepo := userRepository.New(db)
	authentication := service.New(configs.JWTSecret, configs.TokenExp, userRepo)
	userHandler.New(router, authentication, userRepo)

	sessionsRepo := sessionsRepository.New(db, timeZone)
	sessionsServ := sessionsService.New(sessionsRepo)
	sessionsHandler.New(router, sessionsServ)

	log.Fatal(http.ListenAndServe(":"+configs.Port, router))
}
