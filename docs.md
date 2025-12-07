# 1. varietas.go

```
// internal/domain/varietas.go
package domain

import "time"

type VarietasPadi struct {
	ID               int       `json:"id_padi"`
	VarietasKelas    string    `json:"varietas_kelas"`
	Warna            string    `json:"warna"`
	PanjangBijiMM    float64   `json:"panjang_biji_mm"`
	TeksturPermukaan string    `json:"tekstur_permukaan"`
	BentukUjungDaun  string    `json:"bentuk_ujung_daun"`
	WaktuPembuatan   time.Time `json:"waktu_pembuatan"`
}

type VarietasRepository interface {
	Create(data VarietasPadi) (VarietasPadi, error)
	FindByID(id int) (VarietasPadi, error)
	FindAll() ([]VarietasPadi, error)
	Update(data VarietasPadi) (VarietasPadi, error)
	Delete(id int) error
}

type VarietasService interface {
	TambahkanData(data VarietasPadi) (VarietasPadi, error)
	DapatkanDataByID(id int) (VarietasPadi, error)
	DapatkanSemuaData() ([]VarietasPadi, error)
	UbahData(data VarietasPadi) (VarietasPadi, error)
	HapusData(id int) error
}

```