package config

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	ErrNoJWTSecret       = errors.New("missing JWT secret key variable")
	ErrNoPort            = errors.New("missing server port variable")
	ErrNoDataBaseURL     = errors.New("missing database URL variable")
	ErrNoTokenExpiration = errors.New("missing token expiration time variable")
)

type Config struct {
	Port      string
	JWTSecret string
	Db        string
	TokenExp  string
}

func New() Config {
	if err := godotenv.Load(); err != nil {
		log.Print("no .env file found")
	}
	return Config{
		Port:      os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		Db:        os.Getenv("DATABASE_URL"),
		TokenExp:  os.Getenv("TOKEN_EXP_IN_HOURS"),
	}
}

func (c Config) Validate() error {
	if c.Port == "" {
		return ErrNoPort
	}
	if c.JWTSecret == "" {
		return ErrNoJWTSecret
	}
	if c.Db == "" {
		return ErrNoDataBaseURL
	}
	if c.TokenExp == "" {
		return ErrNoTokenExpiration
	}
	return nil
}
