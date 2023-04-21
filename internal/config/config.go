package config

import (
	"os"
)

type Config struct {
	Port string
}

func New() Config {
	return Config{
		Port: getEnv("PORT", "8080"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}
