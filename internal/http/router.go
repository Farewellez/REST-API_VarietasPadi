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

	// ROUTE YANG KAMU MAU:
	mux.HandleFunc("GET /api/varietas", pengamatanHandler.GetAll)

	// Optional: root page biar cantik
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
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

	return mux
}
