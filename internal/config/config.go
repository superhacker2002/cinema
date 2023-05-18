package config

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

var ErrNoEnvFile = errors.New("no .env file found")

type Config struct {
	Port      string `env:"PORT,default=8080"`
	JWTSecret string `env:"JWT_SECRET,default=secret-key"`
	Db        string `env:"DATABASE_URL,default=localhost:5432/cinema"`
	TokenExp  int    `env:"TOKEN_EXP_IN_HOURS,default=24"`
}

func New() (Config, error) {
	var c Config
	if err := godotenv.Load(); err != nil {
		return c, ErrNoEnvFile
	}

	ctx := context.Background()
	if err := envconfig.Process(ctx, &c); err != nil {
		return c, err
	}

	return c, nil
}
