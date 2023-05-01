package main

import (
	cinemaHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinema/handler"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/config"
	userHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/user/handler"
	"log"
	"net/http"
)

func main() {
	config := config.New()
	userHandler.New()
	cinemaHandler.New()
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
