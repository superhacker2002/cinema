package main

import (
	authHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/auth/handler"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/auth/middleware"
	authRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/auth/repository"
	authService "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/auth/service"
	sessionsHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/cinemasession/handler"
	sessionsRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/cinemasession/repository"
	sessionsService "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/cinemasession/service"
	hallsHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/hall/handler"
	hallsRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/hall/repository"
	hallsService "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/hall/service"
	moviesHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/movie/handler"
	moviesRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/movie/repository"
	moviesService "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/movie/service"
	ticketHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/ticket/handler"
	pdfGenerator "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/ticket/pdf"
	ticketRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/ticket/repository"
	ticketService "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/ticket/service"
	userHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/user/handler"
	userRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/user/repository"
	userService "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/user/service"

	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/config"
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

	authRepo := authRepository.New(db)
	authServ := authService.New(configs.JWTSecret, configs.TokenExp, authRepo)
	authHandler.New(router, authServ)

	authorizer := middleware.New(authServ)

	userRepo := userRepository.New(db)
	userServ := userService.New(userRepo)
	userHandler.New(router, userServ)

	sessionsRepo := sessionsRepository.New(db, configs.TimeZone)
	sessionsServ := sessionsService.New(sessionsRepo)
	handler := sessionsHandler.New(sessionsServ)
	handler.SetRoutes(router, authorizer)

	hallsRepo := hallsRepository.New(db)
	hallsServ := hallsService.New(hallsRepo)
	hallsHandler.New(router, hallsServ)

	moviesRepo := moviesRepository.New(db)
	moviesServ := moviesService.New(moviesRepo)
	moviesHandler.New(router, moviesServ)

	ticketGen := pdfGenerator.New()
	ticketRepo := ticketRepository.New(db)
	ticketServ := ticketService.New(ticketRepo, ticketGen)
	tHandler := ticketHandler.New(ticketServ)
	tHandler.SetRoutes(router, authorizer)

	log.Fatal(http.ListenAndServe(":"+configs.Port, router))
}
