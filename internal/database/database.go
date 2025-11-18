// internal/database/postgres.go
package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // v5 — terbaru, cepet, aktif maintained
)

func Connect(dsn string) *sql.DB {
	// "pgx" adalah nama driver untuk database/sql
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to open PostgreSQL connection:", err)
	}

	// Ping biar langsung ketahuan kalau DSN salah
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping PostgreSQL:", err)
	}

	// Tuning connection pool
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Successfully connected to Neon PostgreSQL!")
	return db
}
