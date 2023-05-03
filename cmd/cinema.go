package main

import (
	cinemaHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinema/handler"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/config"
	userHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/user/handler"
	"database/sql"
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

	router := mux.NewRouter()
	userHandler.New(router)
	cinemaHandler.New(router)

	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
