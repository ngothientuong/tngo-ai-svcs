package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Get the absolute path to the .env file
	envPath, err := filepath.Abs("../../.env")
	if err != nil {
		log.Fatalf("Error getting absolute path to .env file: %v", err)
	}

	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
