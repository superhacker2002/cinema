package config

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	ErrNoJWTSecret   = errors.New("missing JWT secret key variable")
	ErrNoPort        = errors.New("missing server port variable")
	ErrNoDataBaseURL = errors.New("missing database URL variable")
)

type Config struct {
	Port string
	Db   string
}

func New() Config {
	if err := godotenv.Load(); err != nil {
		log.Print("no .env file found")
	}
	return Config{
		Port: os.Getenv("PORT"),
		Db:   os.Getenv("DATABASE_URL"),
	}
}

func (c Config) Validate() error {
	if c.Port == "" {
		return ErrNoPort
	}
	if c.Db == "" {
		return ErrNoDataBaseURL
	}
	if c.Db == "" {
		return ErrNoDataBaseURL
	}
	return nil
}
