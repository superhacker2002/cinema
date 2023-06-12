package main

import (
	authHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/auth/handler"
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/auth/middleware"
	authRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/auth/repository"
	authService "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/auth/service"
	sessionsHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/cinemasession/handler"
	sessionsRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/cinemasession/repository"
	sessionsService "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/cinemasession/service"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	hallsHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/hall/handler"
	hallsRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/hall/repository"
	hallsService "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/hall/service"

	moviesHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/movie/handler"
	moviesRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/movie/repository"
	moviesService "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/movie/service"

	ticketHandler "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/ticket/handler"
	pdf "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/ticket/pdf"
	ticketRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/ticket/repository"
	ticketService "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/domains/ticket/service"
	minioStorage "bitbucket.org/Ernst_Dzeravianka/cinemago-app/pkg/minio"

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
		log.Fatalf("failed to open connection with database: %v", err)
	}
	defer db.Close()

	minioClient, err := minio.New(configs.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(configs.MinIOUser, configs.MinIOPasswd, ""),
		Secure: false,
	})

	if err != nil {
		log.Fatalf("failed to create MinIO client: %v", err)
	}

	router := mux.NewRouter()

	authRepo := authRepository.New(db)
	authServ := authService.New(configs.JWTSecret, configs.TokenExp, authRepo)
	authHandler.New(router, authServ)

	authMW := authmw.New(authServ)

	userRepo := userRepository.New(db)
	userServ := userService.New(userRepo)
	userHandler.New(router, userServ)

	sessionsRepo := sessionsRepository.New(db, configs.TimeZone)
	sessionsServ := sessionsService.New(sessionsRepo)
	sessionsHandler.New(sessionsServ).SetRoutes(router, authMW)

	hallsRepo := hallsRepository.New(db)
	hallsServ := hallsService.New(hallsRepo)
	hallsHandler.New(hallsServ).SetRoutes(router, authMW)

	moviesRepo := moviesRepository.New(db)
	moviesServ := moviesService.New(moviesRepo)
	moviesHandler.New(moviesServ).SetRoutes(router, authMW)

	ticketGen := pdf.Generator{}
	ticketsStorage := minioStorage.New(minioClient)

	ticketRepo := ticketRepository.New(db)
	ticketServ := ticketService.New(ticketRepo, ticketGen, ticketsStorage)
	ticketHandler.New(ticketServ).SetRoutes(router, authMW)

	log.Fatal(http.ListenAndServe(":"+configs.Port, router))
}
