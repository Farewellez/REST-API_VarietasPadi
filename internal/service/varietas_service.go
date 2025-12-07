package service

import (
	"context"
	"database/sql" // DITAMBAH: Untuk penanganan error sql.ErrNoRows
	"errors"

	"github.com/Farewellez/REST-API_VarietasPadi/internal/domain"
)

// VarietasService adalah struct implementasi dari domain.VarietasService.
// Struct ini menerima interface Repository sebagai dependency.
type VarietasService struct {
	// Variabel repo harus berupa interface, BUKAN struct konkret
	repo domain.VarietasRepository
}

// NewVarietasService adalah constructor untuk Service Layer.
func NewVarietasService(repo domain.VarietasRepository) domain.VarietasService {
	return &VarietasService{repo: repo}
}

// --- IMPLEMENTASI FUNGSI CRUD LENGKAP ---

// DapatkanSemuaData mengimplementasikan kontrak service untuk Read All.
func (s *VarietasService) DapatkanSemuaData(ctx context.Context) ([]domain.VarietasPadi, error) {
	data, err := s.repo.FindAll(ctx) // DITAMBAH ctx
	if err != nil {
		return nil, errors.New("gagal mengambil data varietas dari penyimpanan")
	}

	// Contoh Logika Bisnis: Memberi notifikasi jika data kosong
	if len(data) == 0 {
		// Log.Println("Tidak ada data varietas ditemukan")
	}

	return data, nil
}

// DapatkanDataByID mengimplementasikan kontrak service untuk Read By ID.
// Didefinisikan di luar struct, sebagai method.
func (s *VarietasService) DapatkanDataByID(ctx context.Context, id int) (domain.VarietasPadi, error) {
	data, err := s.repo.FindByID(ctx, id) // DITAMBAH ctx
	if err != nil {
		// Logika penanganan error database spesifik (contoh)
		if errors.Is(err, sql.ErrNoRows) {
			return domain.VarietasPadi{}, errors.New("data varietas tidak ditemukan")
		}
		return domain.VarietasPadi{}, errors.New("gagal mengambil data dari penyimpanan")
	}
	return data, nil
}

// TambahkanData mengimplementasikan kontrak service untuk Create.
// Tanda tangan fungsi diubah untuk menerima context.Context
func (s *VarietasService) TambahkanData(ctx context.Context, data domain.VarietasPadi) (domain.VarietasPadi, error) { // DITAMBAH ctx
	// Logika Bisnis: Validasi
	if data.VarietasKelas == "" || data.PanjangBijiMM <= 0 {
		return domain.VarietasPadi{}, errors.New("varietas kelas atau panjang biji tidak valid")
	}

	// Panggil Repository (DITAMBAH ctx)
	return s.repo.Create(ctx, data)
}

// UbahData mengimplementasikan kontrak service untuk Update. (Perlu implementasi di sini)
func (s *VarietasService) UbahData(ctx context.Context, data domain.VarietasPadi) (domain.VarietasPadi, error) {
	// Logika Bisnis: Validasi ID dan data
	if data.ID <= 0 {
		return domain.VarietasPadi{}, errors.New("ID varietas tidak valid untuk diubah")
	}
	// Panggil Repository (Update)
	return s.repo.Update(ctx, data)
}

// HapusData mengimplementasikan kontrak service untuk Delete. (Perlu implementasi di sini)
func (s *VarietasService) HapusData(ctx context.Context, id int) error {
	// Panggil Repository (Delete)
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("data varietas yang akan dihapus tidak ditemukan")
	}
	return err
}

// --- IMPLEMENTASI FUNCTIONAL PROGRAMMING (FP) ---

// Tipe Predicate (Fungsi FP untuk kriteria filter)
type VarietasPredicate func(v domain.VarietasPadi) bool

// FilterData adalah helper FP yang dapat digunakan kembali di service.
func FilterData(varietas []domain.VarietasPadi, p VarietasPredicate) []domain.VarietasPadi {
	var filtered []domain.VarietasPadi
	for _, v := range varietas {
		if p(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// DapatkanVarietasBijiPanjang adalah contoh fungsi business logic baru
func (s *VarietasService) DapatkanVarietasBijiPanjang(ctx context.Context) ([]domain.VarietasPadi, error) {
	semuaVarietas, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	isBijiPanjang := func(v domain.VarietasPadi) bool {
		return v.PanjangBijiMM > 7.0
	}

	hasilFilter := FilterData(semuaVarietas, isBijiPanjang)
	return hasilFilter, nil
}
