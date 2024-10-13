package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from the .env file.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
}
