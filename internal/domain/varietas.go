// internal/domain/varietas.go
package domain

import "time"

type PengamatanPadi struct {
	ID               int       `json:"id_padi"`
	VarietasKelas    string    `json:"varietas_kelas"`
	Warna            string    `json:"warna"`
	PanjangBijiMM    float64   `json:"panjang_biji_mm"`
	TeksturPermukaan string    `json:"tekstur_permukaan"`
	BentukUjungDaun  string    `json:"bentuk_ujung_daun"`
	WaktuPembuatan   time.Time `json:"waktu_pembuatan"`
}
