package config

import (
	"os"
)

type Config struct {
	Port        string
	MongoDB_URL string
}

func NewConfig() *Config {
	return &Config{
		Port:        os.Getenv("PORT"),
		MongoDB_URL: os.Getenv("MONGODB_URL"),
	}
}
