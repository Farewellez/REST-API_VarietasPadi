// internal/domain/varietas.go
package domain

import (
	"context" // WAJIB: Import context karena digunakan di Interface
	"time"
)

type VarietasPadi struct {
	ID               int       `json:"id_padi"`
	VarietasKelas    string    `json:"varietas_kelas"`
	Warna            string    `json:"warna"`
	PanjangBijiMM    float64   `json:"panjang_biji_mm"`
	TeksturPermukaan string    `json:"tekstur_permukaan"`
	BentukUjungDaun  string    `json:"bentuk_ujung_daun"`
	WaktuPembuatan   time.Time `json:"waktu_pembuatan"`
}

// VarietasRepository Interface (Kontrak Data Access)
type VarietasRepository interface {
	// SEMUA FUNGSI CRUD DITAMBAH context.Context SEBAGAI ARGUMEN PERTAMA
	Create(ctx context.Context, data VarietasPadi) (VarietasPadi, error)
	FindByID(ctx context.Context, id int) (VarietasPadi, error)
	FindAll(ctx context.Context) ([]VarietasPadi, error) // FIX ERROR: Menambah context.Context
	Update(ctx context.Context, data VarietasPadi) (VarietasPadi, error)
	Delete(ctx context.Context, id int) error
}

// VarietasService Interface (Kontrak Logika Bisnis)
type VarietasService interface {
	// SEMUA FUNGSI SERVICE DITAMBAH context.Context SEBAGAI ARGUMEN PERTAMA
	// Ini adalah kontrak lengkap untuk CRUD (sudah benar, hanya perlu context)
	TambahkanData(ctx context.Context, data VarietasPadi) (VarietasPadi, error)
	DapatkanDataByID(ctx context.Context, id int) (VarietasPadi, error)    // FIX ERROR: Menambah context.Context
	DapatkanSemuaData(ctx context.Context) ([]VarietasPadi, error)         // FIX ERROR: Menambah context.Context
	UbahData(ctx context.Context, data VarietasPadi) (VarietasPadi, error) // FIX ERROR: Menambah context.Context
	HapusData(ctx context.Context, id int) error                           // FIX ERROR: Menambah context.Context
}
