package config

import (
	"errors" // Import untuk mengembalikan error
	"os"

	"github.com/joho/godotenv"
)

// Config adalah struct konfigurasi aplikasi
type Config struct {
	DBURL string // connection string ke Neon PostgreSQL
	Port  string // port server (default 8080)
}

// Load membaca konfigurasi dari environment variable atau .env
// Mengembalikan Config dan error jika gagal memuat konfigurasi wajib.
func Load() (Config, error) { // Mengubah Load agar mengembalikan (Config, error)
	// Load file .env kalau ada (nggak error kalau nggak ada)
	// Kita abaikan error godotenv.Load() karena tidak wajib ada
	_ = godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		// Mengembalikan error, membiarkan main.go yang handle log.Fatal
		return Config{}, errors.New("DB_URL tidak ditemukan! Harap set environment variable atau file .env")
	}

	// Kalau mau nambah port custom, bisa dari env juga
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default
	}

	return Config{
		DBURL: dbURL,
		Port:  port,
	}, nil // Mengembalikan nil (tidak ada error)
}
