package main

import (
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
	config := config.New()
	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", config.Db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	userRepository := userRepository.New(db)
	cinemaRepository := cinemaRepository.New(db)

	router := mux.NewRouter()

	userHandler.New(router, userRepository)
	cinemaHandler.New(router, cinemaRepository)

	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
