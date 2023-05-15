package main

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/auth"
	cinemaHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinema/handler"
	cinemaRepo "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinema/repository"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/config"

	userHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/user/handler"
	userRepo "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/user/repository"
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func main() {
	log.SetFlags(log.Lshortfile)
	config := config.New()
	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", config.Db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	cinemaRepo := cinemaRepo.New(db)
	userRepo := userRepo.New(db)

	router := mux.NewRouter()
	exp, err := strconv.Atoi(config.TokenExp)
	if err != nil {
		log.Fatal(err)
	}
	authentication := auth.New(config.JWTSecret, exp, userRepo)

	userHandler.New(router, authentication, userRepo)
	cinemaHandler.New(router, cinemaRepo)

	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
