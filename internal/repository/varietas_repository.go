// internal/repository/pengamatan_repository.go
package repository

import (
	"database/sql"

	"github.com/Farewellez/REST-API_VarietasPadi/internal/domain"
)

type PengamatanRepository struct {
	db *sql.DB
}

func NewPengamatanRepository(db *sql.DB) *PengamatanRepository {
	return &PengamatanRepository{db: db}
}

func (r *PengamatanRepository) FindAll() ([]domain.PengamatanPadi, error) {
	query := `
		SELECT id_padi, varietas_kelas, warna, panjang_biji_mm, 
		       tekstur_permukaan, bentuk_ujung_daun, waktu_pembuatan
		FROM DataPengamatanPadi
		ORDER BY id_padi
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.PengamatanPadi
	for rows.Next() {
		var p domain.PengamatanPadi
		if err := rows.Scan(&p.ID, &p.VarietasKelas, &p.Warna, &p.PanjangBijiMM,
			&p.TeksturPermukaan, &p.BentukUjungDaun, &p.WaktuPembuatan); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, rows.Err()
}
