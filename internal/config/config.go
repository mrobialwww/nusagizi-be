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

	// Auth0 M2M credentials (untuk Management API)
	Auth0ClientID     string
	Auth0ClientSecret string

	// Auth0 Role IDs
	Auth0RoleIDMother    string
	Auth0RoleIDCaregiver string
	Auth0RoleIDDoctor    string
}

func Load() (*Config, error) {
	var err error = godotenv.Load()

	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	var config *Config = &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		Port:          os.Getenv("PORT"),
		Auth0Domain:   os.Getenv("AUTH0_DOMAIN"),
		Auth0Audience: os.Getenv("AUTH0_AUDIENCE"),

		Auth0ClientID:     os.Getenv("CLIENT_ID"),
		Auth0ClientSecret: os.Getenv("CLIENT_SECRET"),

		Auth0RoleIDMother:    os.Getenv("AUTH0_ROLE_ID_MOTHER"),
		Auth0RoleIDCaregiver: os.Getenv("AUTH0_ROLE_ID_CAREGIVER"),
		Auth0RoleIDDoctor:    os.Getenv("AUTH0_ROLE_ID_DOCTOR"),
	}

	return config, nil
}