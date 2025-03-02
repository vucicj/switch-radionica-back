package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	JWTSecret            string
	TokenDuration        time.Duration
	RefreshTokenDuration time.Duration
}

func LoadConfig() *Config {
	err := godotenv.Load() // Load .env file
	if err != nil {
		log.Println("No .env file found, using system env vars")
	}

	return &Config{
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", "postgres"),
		DBPassword:           getEnv("DB_PASSWORD", "yourpassword"),
		DBName:               getEnv("DB_NAME", "mydb"),
		JWTSecret:            getEnv("JWT_SECRET", "your-secret-key"),
		TokenDuration:        getEnvDuration("TOKEN_DURATION", 15*time.Minute),
		RefreshTokenDuration: getEnvDuration("REFRESH_TOKEN_DURATION", 7*24*time.Hour),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if dur, err := time.ParseDuration(value); err == nil {
			return dur
		}
		log.Printf("Invalid duration for %s: %s, using fallback", key, value)
	}
	return fallback
}
