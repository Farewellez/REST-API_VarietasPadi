// internal/repository/varietas_repository.go
package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Farewellez/REST-API_VarietasPadi/internal/domain"
)

type VarietasRepository struct {
	db *sql.DB
}

func NewVarietasRepository(db *sql.DB) *VarietasRepository {
	return &VarietasRepository{db: db}
}

func (r *VarietasRepository) FindAll(ctx context.Context) ([]domain.VarietasPadi, error) {
	query := `
		SELECT id_padi, varietas_kelas, warna, panjang_biji_mm, 
		       tekstur_permukaan, bentuk_ujung_daun, waktu_pembuatan
		FROM DataPengamatanPadi
		ORDER BY id_padi
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.VarietasPadi
	for rows.Next() {
		var p domain.VarietasPadi
		if err := rows.Scan(&p.ID, &p.VarietasKelas, &p.Warna, &p.PanjangBijiMM,
			&p.TeksturPermukaan, &p.BentukUjungDaun, &p.WaktuPembuatan); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, rows.Err()
}

// Mengimplementasikan interface domain.VarietasRepository
func (r *VarietasRepository) Create(ctx context.Context, data domain.VarietasPadi) (domain.VarietasPadi, error) {
	query := `
        INSERT INTO DataPengamatanPadi (varietas_kelas, warna, panjang_biji_mm, 
                                      tekstur_permukaan, bentuk_ujung_daun)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id_padi
    `
	err := r.db.QueryRow(query,
		data.VarietasKelas,
		data.Warna,
		data.PanjangBijiMM,
		data.TeksturPermukaan,
		data.BentukUjungDaun,
	).Scan(&data.ID)

	if err != nil {
		return domain.VarietasPadi{}, err
	}
	return data, nil
}

// internal/repository/varietas_repository.go (Tambahan)

func (r *VarietasRepository) FindByID(ctx context.Context, id int) (domain.VarietasPadi, error) {
	query := `
        SELECT id_padi, varietas_kelas, warna, panjang_biji_mm, 
               tekstur_permukaan, bentuk_ujung_daun, waktu_pembuatan
        FROM DataPengamatanPadi
        WHERE id_padi = $1
    `
	var p domain.VarietasPadi

	// Gunakan QueryRowContext untuk operasi Read tunggal
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.VarietasKelas, &p.Warna, &p.PanjangBijiMM,
		&p.TeksturPermukaan, &p.BentukUjungDaun, &p.WaktuPembuatan,
	)

	if err != nil {
		// Jika data tidak ditemukan, kembalikan error spesifik dari sql
		if errors.Is(err, sql.ErrNoRows) {
			return domain.VarietasPadi{}, sql.ErrNoRows
		}
		return domain.VarietasPadi{}, err
	}
	return p, nil
}

// internal/repository/varietas_repository.go (Tambahan)

func (r *VarietasRepository) Update(ctx context.Context, data domain.VarietasPadi) (domain.VarietasPadi, error) {
	query := `
        UPDATE DataPengamatanPadi
        SET varietas_kelas=$2, warna=$3, panjang_biji_mm=$4, 
            tekstur_permukaan=$5, bentuk_ujung_daun=$6
        WHERE id_padi = $1
        RETURNING id_padi, waktu_pembuatan
    `

	// Gunakan QueryRowContext untuk mendapatkan kembali data yang baru diupdate
	err := r.db.QueryRowContext(ctx, query,
		data.ID,
		data.VarietasKelas,
		data.Warna,
		data.PanjangBijiMM,
		data.TeksturPermukaan,
		data.BentukUjungDaun,
	).Scan(&data.ID, &data.WaktuPembuatan)

	if err != nil {
		return domain.VarietasPadi{}, err
	}
	return data, nil
}

// internal/repository/varietas_repository.go (Tambahan)

func (r *VarietasRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM DataPengamatanPadi WHERE id_padi = $1`

	// Gunakan ExecContext untuk operasi yang tidak mengembalikan row
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Cek apakah ada row yang terpengaruh (terhapus)
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Ini penting: jika tidak ada row yang terhapus, kita kembalikan sql.ErrNoRows
		return sql.ErrNoRows
	}

	return nil
}
