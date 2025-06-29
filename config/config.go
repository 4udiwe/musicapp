package config

import (
	"os"
)

type Config struct {
	PostgresURL string
	ServerPort  string
}

func LoadConfig() *Config {
	return &Config{
		PostgresURL: os.Getenv("PG_URL"),
		ServerPort:  os.Getenv("SERVER_PORT"),
	}
}
