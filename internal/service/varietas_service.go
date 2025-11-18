// internal/http/handler/pengamatan_handler.go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Farewellez/REST-API_VarietasPadi/internal/repository"
)

type PengamatanHandler struct {
	repo *repository.PengamatanRepository
}

func NewPengamatanHandler(repo *repository.PengamatanRepository) *PengamatanHandler {
	return &PengamatanHandler{repo: repo}
}

func (h *PengamatanHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.repo.FindAll()
	if err != nil {
		http.Error(w, "Gagal mengambil data dari database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"total":   len(data),
		"data":    data,
	})
}
