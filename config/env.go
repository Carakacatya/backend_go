package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv digunakan untuk memuat konfigurasi dari file .env
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}
}
