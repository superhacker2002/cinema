package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

var (
	ErrNoEnvFile          = errors.New("no .env file found")
	ErrNoJWTSecret        = errors.New("missing JWT secret key variable")
	ErrNoPort             = errors.New("missing server port variable")
	ErrNoDataBaseURL      = errors.New("missing database URL variable")
	ErrBadTokenExpiration = errors.New("missing or incorrect token expiration time variable")
)

type Config struct {
	Port      string
	JWTSecret string
	Db        string
	TokenExp  int
}

func New() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, ErrNoEnvFile
	}
	config := Config{
		Port:      os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		Db:        os.Getenv("DATABASE_URL"),
	}
	tokenExp, err := getEnvInt("TOKEN_EXP_IN_HOURS")
	if err != nil {
		return config, ErrBadTokenExpiration
	}
	config.TokenExp = tokenExp
	return config, nil
}

func getEnvInt(varName string) (int, error) {
	varString := os.Getenv(varName)
	varInt, err := strconv.Atoi(varString)
	if err != nil {
		return 0, err
	}
	return varInt, nil
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
	if c.TokenExp <= 0 {
		return ErrBadTokenExpiration
	}
	return nil
}
