package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 int
	DatabaseURL          string
	GoogleClientID       string
	GoogleClientSecret   string
	SessionSecret        string
	UploadsDir          string
	Environment         string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		port = 8080
	}

	return &Config{
		Port:                port,
		DatabaseURL:         getEnv("DATABASE_URL", "postgres://tipjar:tipjar@localhost/tipjar?sslmode=disable"),
		GoogleClientID:      os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:  os.Getenv("GOOGLE_CLIENT_SECRET"),
		SessionSecret:       getEnv("SESSION_SECRET", "your-secret-key-change-this"),
		UploadsDir:         getEnv("UPLOADS_DIR", "./uploads"),
		Environment:        getEnv("ENVIRONMENT", "development"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
