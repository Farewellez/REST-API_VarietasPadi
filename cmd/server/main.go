package main

import (
	"context" // Tambahkan context
	"log"
	"net/http"
	"time"

	"github.com/Farewellez/REST-API_VarietasPadi/internal/config"
	"github.com/Farewellez/REST-API_VarietasPadi/internal/database"
	httpHandler "github.com/Farewellez/REST-API_VarietasPadi/internal/http"
	"github.com/Farewellez/REST-API_VarietasPadi/internal/http/handler"
	"github.com/Farewellez/REST-API_VarietasPadi/internal/repository"
	"github.com/Farewellez/REST-API_VarietasPadi/internal/service"
)

func main() {
	// 1. MEMUAT KONFIGURASI
	// Menggunakan pola error handling yang benar dari config.Load()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("FATAL: Gagal memuat konfigurasi: ", err)
	}

	// 2. KONEKSI KE DATABASE (PostgreSQL/NeonDB)
	// Menggunakan pola error handling yang benar dari database.NewDB()
	db, err := database.NewDB(cfg.DBURL)
	if err != nil {
		log.Fatal("FATAL: Gagal terhubung ke database: ", err)
	}
	defer db.Close()
	log.Println("Berhasil terhubung ke database!")

	// Opsional: Cek lagi status DB sebelum start
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("FATAL: Database ping gagal setelah koneksi: ", err)
	}

	// 3. WIRING UP (Inisialisasi Lapisan)
	// A. Inisialisasi Repository
	varietasRepo := repository.NewVarietasRepository(db)

	// B. Inisialisasi Service (DI: Membutuhkan Repository Interface)
	varietasService := service.NewVarietasService(varietasRepo)

	// C. Inisialisasi Handler (DI: Membutuhkan Service Interface)
	varietasHandler := handler.NewVarietasHandler(varietasService)

	// D. Inisialisasi Router (Membutuhkan Handler)
	router := httpHandler.NewRouter(varietasHandler)

	// 4. MENJALANKAN SERVER
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Server siap dijalankan pada port %s (http://localhost%s)", cfg.Port, srv.Addr)

	// Blok dan tunggu traffic masuk
	log.Fatal(srv.ListenAndServe())
}
