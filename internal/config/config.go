package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    DatabaseURL   string
    Port          string
    Auth0Domain   string
    Auth0Audience string
}

func Load() (*Config, error) {
	var err error = godotenv.Load()

	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	var config *Config = &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
		Auth0Domain:   os.Getenv("AUTH0_DOMAIN"),
  Auth0Audience: os.Getenv("AUTH0_AUDIENCE"),
	}

	return config, nil
}