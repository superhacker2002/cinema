package config

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	ErrNoJWTSecret = errors.New("missing JWT secret key variable")
	ErrNoPort      = errors.New("missing server port variable")
)

type Config struct {
	Port      string
	JWTSecret []byte
}

func New() Config {
	if err := godotenv.Load(); err != nil {
		log.Print("no .env file found")
	}
	return Config{
		Port:      os.Getenv("PORT"),
		JWTSecret: []byte(os.Getenv("JWT_SECRET")),
	}
}

func (c Config) Validate() error {
	if c.Port == "" {
		return ErrNoPort
	}
	if len(c.JWTSecret) == 0 {
		return ErrNoJWTSecret
	}
	return nil
}
