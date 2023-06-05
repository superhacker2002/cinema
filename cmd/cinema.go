package main

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/auth/service"
	sessionsHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasession/handler"
	sessionsRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasession/repository"
	sessionsService "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasession/service"

	hallsHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/hall/handler"
	hallsRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/hall/repository"
	hallsService "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/hall/service"

	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/config"
	userHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/user/handler"
	userRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/user/repository"

	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Lshortfile)
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

	sessionsRepo := sessionsRepository.New(db, configs.TimeZone)
	sessionsServ := sessionsService.New(sessionsRepo)
	sessionsHandler.New(router, sessionsServ)

	hallsRepo := hallsRepository.New(db)
	hallsServ := hallsService.New(hallsRepo)
	hallsHandler.New(router, hallsServ)

	log.Fatal(http.ListenAndServe(":"+configs.Port, router))
}
