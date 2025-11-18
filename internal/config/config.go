// internal/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv" // biar bisa baca file .env (opsional tapi super enak)
)

// Config adalah struct konfigurasi aplikasi
type Config struct {
	DBURL string // connection string ke Neon PostgreSQL
	Port  string // port server (default :8080)
}

// Load membaca konfigurasi dari environment variable atau .env
func Load() Config {
	// Load file .env kalau ada (nggak error kalau nggak ada)
	_ = godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL tidak ditemukan! Set dulu di environment atau file .env")
	}

	// Kalau mau nambah port custom, bisa dari env juga
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default
	}

	return Config{
		DBURL: dbURL,
		Port:  port,
	}
}
