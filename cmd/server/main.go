package main

import (
	"log"
	"net/http"
	"time"

	/*
		Jadi penggunaan github.com/Farewellez/REST-API_VarietasPadi disini karena mengikuti go mod init di awal. Jadi ketika ingin
		import path tertentu maka perlu mencantumkan full modul path di awal initialize.

		Untuk bagian hhtphandler disini sebagai alias di Go untuk paket/path yang di import. Karena ada dua path http yang sama disini.
		net/http (modul GO) dan juga internal/http (path projek)
	*/
	"github.com/Farewellez/REST-API_VarietasPadi/internal/config"
	"github.com/Farewellez/REST-API_VarietasPadi/internal/database"
	httpHandler "github.com/Farewellez/REST-API_VarietasPadi/internal/http"
)

func main() {
	/*
		Step1: memanggil function Load yang ada di config.go. Load disini akan mengembalikkan type struct Config yang memiliki 2 isi field yaitu
		DBUrl dan Port yang akan disimpan di variable cfg
	*/
	cfg := config.Load() // buat load configurasi dari env config di internal
	/*
		Step2: setelah informasi port dan DBUrl didapat maka bisa melanjutkan ke koneksi database menggunakan function connect() dari path database
	*/
	db := database.Connect(cfg.DBURL) // mirip config bedanya cuma di direktori database. bagian ini buat connect ke db, return object maupun handle connection
	/*
		Step3: defer disini adalah statement yang menjadwalkan pemanggilan fungsi dieksekusi segera sebelum return function yang mana artinya ketika aplikasi sudah
		berhenti running maka koneksi ke database baru ditutup
	*/
	defer db.Close() // clean-up kalau udah selesai
	/*
		Step4: defer disini adalah statement yang menjadwalkan pemanggilan fungsi dieksekusi segera sebelum return function yang mana artinya ketika aplikasi sudah
		berhenti running maka koneksi ke database baru ditutup
	*/
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
