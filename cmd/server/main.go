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
	cfg := config.Load()
	db := database.Connect(cfg.DBURL)
	defer db.Close()

	router := httpHandler.NewRouter(db)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server running on port: %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
