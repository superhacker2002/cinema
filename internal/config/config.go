package config

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"time"
)

var timeZone = time.FixedZone("UTC+4", 4*60*60)

type Config struct {
	Port      string `env:"PORT,default=8080"`
	JWTSecret string `env:"JWT_SECRET,default=secret-key"`
	Db        string `env:"DATABASE_URL,default=postgresql://postgres:2587@localhost:5432/cinema?sslmode=disable"`
	TokenExp  int    `env:"TOKEN_EXP_IN_HOURS,default=24"`
	TimeZone  *time.Location
}

func New() (Config, error) {
	var c Config
	if err := godotenv.Load(); err != nil {
		return c, err
	}

	ctx := context.Background()
	if err := envconfig.Process(ctx, &c); err != nil {
		return c, err
	}

	c.TimeZone = timeZone

	return c, nil
}
