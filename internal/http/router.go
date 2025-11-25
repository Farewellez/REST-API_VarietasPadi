// internal/http/router.go
package http

import (
	"database/sql"
	"net/http"

	"github.com/Farewellez/REST-API_VarietasPadi/internal/http/handler"
	"github.com/Farewellez/REST-API_VarietasPadi/internal/repository"
)

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// Repository & Handler
	pengamatanRepo := repository.NewPengamatanRepository(db)
	pengamatanHandler := handler.NewPengamatanHandler(pengamatanRepo)

	// ONLY endpoint that exists
	mux.HandleFunc("/api/varietas", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		pengamatanHandler.GetAll(w, r)
	})

	// Serve index.html
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
		  "message": "REST API Varietas & Pengamatan Padi",
		  "status": "active",
		  "endpoints": {
		    "GET": "/api/varietas → ambil semua data pengamatan padi"
		  },
		  "docs": "https://github.com/Farewellez/REST-API_VarietasPadi"
		}`))
	})

	// Serve static files
	fs := http.FileServer(http.Dir("views/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return mux
}
