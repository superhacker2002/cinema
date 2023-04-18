package main

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/config"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/controller/httphandler"
	"log"
	"net/http"
)

func main() {
	config := config.New()
	httphandler.New()
	log.Fatal(http.ListenAndServe("localhost:"+config.Cinema.Port, nil))
}
