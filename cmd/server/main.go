package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Farewellez/REST-API_VarietasPadi/internal/config"
	"github.com/Farewellez/REST-API_VarietasPadi/internal/database"
	httpHandler "github.com/Farewellez/REST-API_VarietasPadi/internal/http"
)

func main() {
	cfg := config.Load()                // buat load configurasi dari env config di internal
	db := database.Connect(cfg.DBURL)   // mirip config bedanya cuma di direktori database. bagian ini buat connect ke db, return object maupun handle connection
	defer db.Close()                    // clean-up kalau udah selesai
	router := httpHandler.NewRouter(db) // ini buat ngambil beberapa object dari directori http. Ini juga buat ngambil semua route yang dibuat

	// srv bakal pointer ke alamat http.server
	srv := &http.Server{
		Addr:         ":" + cfg.Port, // sekarang ambil dari config
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server running on port: %s", srv.Addr) // cuma buat debug log ajah
	log.Fatal(srv.ListenAndServe())                    // ListenAndServe memiliki side effect (I/O) dan blocking
}
