package main

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/auth"
	cinemaHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinema/handler"
	cinemaRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinema/repository"
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
	userRepo := userRepository.New(db)
	cinemaRepo := cinemaRepository.New(db)

	router := mux.NewRouter()
	authentication := auth.New(configs.JWTSecret, configs.TokenExp, userRepo)

	userHandler.New(router, authentication, userRepo)
	cinemaHandler.New(router, cinemaRepo)

	log.Fatal(http.ListenAndServe(":"+configs.Port, router))
}
