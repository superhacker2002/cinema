package config

import (
	"os"
)

type Config struct {
	Port      string
	JWTSecret []byte
}

func New() Config {
	return Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: []byte(os.Getenv("JWT_SECRET")),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}
