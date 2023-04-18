package config

import (
	"os"
)

type CinemaConfig struct {
	Port string
}

type Config struct {
	Cinema CinemaConfig
}

func New() *Config {
	return &Config{
		Cinema: CinemaConfig{
			Port: getEnv("PORT", "8080"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
