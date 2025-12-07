package http

import (
	"net/http"

	"github.com/gorilla/mux"

	// Import Handler yang sudah kita buat sebelumnya
	"github.com/Farewellez/REST-API_VarietasPadi/internal/http/handler"
)

// NewRouter membuat dan menginisialisasi rute-rute aplikasi
func NewRouter(varietasHandler *handler.VarietasHandler) *mux.Router {
	r := mux.NewRouter()

	// 1. PENANGANAN ASSET STATIS (CSS, JS, GAMBAR)
	// Menyajikan semua file di dalam direktori 'views/'
	// Jika file static ada di 'views/css/style.css', akan diakses melalui /static/css/style.css
	fs := http.FileServer(http.Dir("views"))

	// Kita gunakan Handle di sini agar mux bisa menyajikan file dari direktori
	// http.StripPrefix memastikan path yang diakses tidak memiliki '/static/' di dalamnya saat mencari file di folder 'views'
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Asumsi: Kita asumsikan semua aset (CSS, JS) yang direferensikan oleh index.html
	// berada di dalam folder 'views' atau subfolder 'views/static'.
	// Jika aset kamu ada di 'views/static', ubah baris di atas menjadi:
	// fs := http.FileServer(http.Dir("views/static"))
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// 2. PENANGANAN ROOT (/) -> Menyajikan index.html
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Cek path, jika bukan root (misalnya '/favicon.ico'), biarkan handler lain atau mux handle 404
		if r.URL.Path != "/" {
			// Kita biarkan logika ini di sini, tapi pastikan tidak bentrok dengan /static/
			// Untuk index.html kita hanya ingin melayani path tepat "/"
			return
		}
		// Set Header dan kirimkan file index.html
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "views/index.html")
	})

	// 3. PENANGANAN API
	// Sub-router untuk semua endpoint API Varietas
	api := r.PathPrefix("/api/varietas").Subrouter()

	// --- Rute CRUD Varietas Padi ---
	api.HandleFunc("", varietasHandler.GetAll).Methods(http.MethodGet)
	api.HandleFunc("", varietasHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/{id}", varietasHandler.GetByID).Methods(http.MethodGet)
	api.HandleFunc("/{id}", varietasHandler.Update).Methods(http.MethodPut)
	api.HandleFunc("/{id}", varietasHandler.Delete).Methods(http.MethodDelete)

	return r
}
