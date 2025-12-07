package database

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewDB membuat dan menguji koneksi ke PostgreSQL
func NewDB(dsn string) (*sql.DB, error) {
	// 1. sql.Open
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, errors.New("gagal membuka koneksi PostgreSQL: " + err.Error())
	}

	// 2. Ping (Uji Koneksi Awal)
	if err = db.Ping(); err != nil {
		// Tutup koneksi yang setengah terbuka jika ping gagal
		db.Close()
		return nil, errors.New("gagal melakukan ping ke PostgreSQL: " + err.Error())
	}

	// 3. Tuning Connection Pool
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5) // Biasanya MaxIdle dibuat lebih kecil dari MaxOpen
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
