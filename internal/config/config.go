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
	// To Validate claim "aud" in ID Token from app
	Auth0ClientID     string
	Auth0ClientSecret string

	// Auth0 Role IDs
	Auth0RoleIDMother    string
	Auth0RoleIDCaregiver string
	Auth0RoleIDDoctor    string

	// AI Service
	AIServiceBaseURL string

	// To Communicate with Auth0 APIs
	M2MClientID     string
	M2MClientSecret string

	// Cloudflare R2
	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2Bucket          string
	R2PublicURL       string
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

		AIServiceBaseURL: os.Getenv("AI_SERVICE_BASE_URL"),

		M2MClientID:     os.Getenv("AUTH0_M2M_CLIENT_ID"),
		M2MClientSecret: os.Getenv("AUTH0_M2M_CLIENT_SECRET"),

		// Cloudflare R2
		R2AccountID:       os.Getenv("R2_ACCOUNT_ID"),
		R2AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		R2SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
		R2Bucket:          os.Getenv("R2_BUCKET"),
		R2PublicURL:       os.Getenv("R2_PUBLIC_URL"),
	}

	return config, nil
}
