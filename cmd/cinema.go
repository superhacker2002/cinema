package main

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/controller/httphandler"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("no .env file found")
	}
}

func main() {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Fatal("no server port provided in .env file")
	}

	httphandler.New()
	log.Fatal(http.ListenAndServe("localhost:"+port, nil))
}
